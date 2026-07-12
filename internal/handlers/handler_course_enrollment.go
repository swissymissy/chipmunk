package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
	"github.com/swissymissy/chipmunk/internal/middleware"
)

// let student enroll for courses
func (cfg *ApiConfig) HandlerEnrollment(w http.ResponseWriter, r *http.Request) {
	// check for studentID
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// decode enroll request
	var req NewEnrollmentRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding new enrollment request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to enroll to the course")
		return
	}

	// add student to course
	_, err = cfg.DB.NewEnrollment(r.Context(), database.NewEnrollmentParams{
		StudentID: studentID,
		CourseID:  req.CourseID,
	})
	if err != nil {
		log.Printf("error adding student to course: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
