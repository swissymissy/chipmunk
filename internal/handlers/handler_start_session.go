package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/swissymissy/chipmunk/internal/auth"
	"github.com/swissymissy/chipmunk/internal/database"
)

// professor starts a new session for a course
func (cfg *ApiConfig) HandlerStartSession(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req StartSessionRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding start session request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// check if an active session for this course already exists
	_, err = cfg.DB.GetActiveSession(r.Context(), req.CourseID)
	if err == nil {
		log.Printf("attempt to create new session for an already exist session of course %s\n", req.CourseID)
		ResponseWithError(w, http.StatusBadRequest, "already exists session for this course")
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Printf("error checking for active sessions :%s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// generate a secret key for this new session
	secretKey, err := auth.GenSecretKey()
	if err != nil {
		log.Printf("error generate new key: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// create new session for the course
	session, err := cfg.DB.CreateSession(r.Context(), database.CreateSessionParams{
		CourseID:     req.CourseID,
		SessionDate:  time.Now().Format("2006-01-02"),
		SecretKey:    secretKey,
		ClassroomLat: ToNullFloat(req.ClassroomLat),
		ClassroomLng: ToNullFloat(req.ClassroomLng),
		RadiusMeters: sql.NullInt64{Int64: 100, Valid: true},
	})
	if err != nil {
		log.Printf("error creating new session: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to create new session")
		return
	}

	// pre-populate the records
	err = cfg.DB.CreateRecords(r.Context(), database.CreateRecordsParams{
		SessionID: session.ID,
		CourseID:  session.CourseID,
	})
	if err != nil {
		log.Printf("error creating records for new session: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to create records for new session")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusCreated, StartSessionResponse{
		SessionID:   session.ID,
		CourseID:    session.CourseID,
		SessionDate: session.SessionDate,
		Status:      session.Status,
		StartedAt:   LocalizeSQLiteTime(session.StartedAt),
	})
}
