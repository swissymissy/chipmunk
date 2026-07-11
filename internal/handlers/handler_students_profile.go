/*
Show a student's profile page
*/

package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/middleware"
)

func (cfg *ApiConfig) HandlerGetStudentProfile(w http.ResponseWriter, r *http.Request) {
	// check student's ID
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// retrieve student's profile from database
	studentProfile, err := cfg.DB.GetProfileByID(r.Context(), studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("student does not exist: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "unable to find student profile")
			return
		}
		log.Printf("error getting student profile from db: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to get student profile")
		return
	}

	ResponseWithJSON(w, http.StatusOK, StudentProfileResponse{
		SchoolID:  studentProfile.StudentID,
		FirstName: studentProfile.FirstName,
		LastName:  studentProfile.LastName,
		Email:     studentProfile.Email,
		Specialty: studentProfile.Specialty.String,
	})
}
