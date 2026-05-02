package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/auth"
	"github.com/swissymissy/chipmunk/internal/database"
	"github.com/swissymissy/chipmunk/internal/middleware"
)

func (cfg *ApiConfig) HandlerStudentCheckIn(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req StudentCheckinReq
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding student checkin request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	qrToken := req.QRToken

	// extract sessionID and timestamp from token
	payload, _, err := auth.SplitToken(qrToken)
	if err != nil {
		log.Printf("error splitting qr token into parts: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid token")
		return
	}
	sessionID, _, err := auth.ExtractSessionIDTimestamp(string(payload))
	if err != nil {
		log.Printf("error extracting session id and timestamp from qr token: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid token")
		return
	}

	// get sessionID secret and check if session is active
	session, err := cfg.DB.GetSessionByID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("attempt to get non-exist session: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "session not found")
			return
		}
		log.Printf("error fetching session by ID: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to fetch session")
		return
	}
	if session.Status != "active" {
		ResponseWithError(w, http.StatusBadRequest, "incorrect session")
		return
	}

	// validate qr token
	valid := auth.ValidateQRToken(qrToken, session.SecretKey)
	if !valid {
		ResponseWithError(w, http.StatusUnauthorized, "invalid token or token expired")
		return
	}

	// get student ID
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// check if student enroll in this course (only allows enrolled students)
	enrolled, err := cfg.DB.IsEnrolled(r.Context(), database.IsEnrolledParams{
		StudentID: studentID,
		CourseID:  session.CourseID,
	})
	if err != nil {
		log.Printf("error checking if student enrolls in a course: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if enrolled != 1 {
		ResponseWithError(w, http.StatusUnauthorized, "you are not enrolled in this course")
		return
	}

	// calculate the distance between student's coord and classroom's coord
	sP := Point{
		Lat: req.StudentLat,
		Lng: req.StudentLng,
	}
	cP := Point{
		Lat: session.ClassroomLat.Float64,
		Lng: session.ClassroomLng.Float64,
	}
	distance := Haversine(cP, sP)
	if distance > float64(session.RadiusMeters.Int64) {
		ResponseWithError(w, http.StatusBadRequest, "student is too far away from class")
		return
	}

	// record attendance
	checkin, err := cfg.DB.StudentCheckIn(r.Context(), database.StudentCheckInParams{
		StudentLat: ToNullFloat(req.StudentLat),
		StudentLng: ToNullFloat(req.StudentLng),
		Accuracy:   ToNullFloat(req.Accuracy),
		SessionID:  session.ID,
		StudentID:  studentID,
	})
	if err != nil {
		log.Printf("error updating student attendance record: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to check-in")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, StudentCheckInRep{
		Status:    checkin.Status,
		CheckInAt: checkin.CheckInAt.String,
	})
}
