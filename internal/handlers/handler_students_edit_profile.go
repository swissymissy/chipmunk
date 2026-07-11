// let student edit their information
// each query has an endpoint for it
// therefore there are 4 handlers for it
// 1. HandlerStudentUpdateEmail
// 2. HandlerStudentUpdateSchoolID
// 3. HandlerStudentUpdateName
// 4. HandlerStudentRemoveACourse

package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/swissymissy/chipmunk/internal/database"
	"github.com/swissymissy/chipmunk/internal/middleware"
)

// let student update their school ID
func (cfg *ApiConfig) HandlerStudentUpdateSchoolID(w http.ResponseWriter, r *http.Request) {
	// get student ID from the context
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// decode for edit school ID
	var newSchoolID UpdateSchoolIDRequest
	err := DecodeRequest(r, &newSchoolID)
	if err != nil {
		log.Printf("error decoding update school ID request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to update school ID")
		return
	}

	// update student school ID
	studentProfile, err := cfg.DB.UpdateStudentSchoolID(r.Context(), database.UpdateStudentSchoolIDParams{StudentID: newSchoolID.SchoolID, ID: studentID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("student not found: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "unable to update school ID")
			return
		}
		log.Printf("error updating student school ID: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to update student school ID")
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

// let student edit/update their email
func (cfg *ApiConfig) HandlerStudentUpdateEmail(w http.ResponseWriter, r *http.Request) {
	// get student ID from the context
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// decode edit email request
	var newEmail UpdateEmailRequest
	err := DecodeRequest(r, &newEmail)
	if err != nil {
		log.Printf("error decoding update email request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to update email")
		return
	}

	// check student's input
	email, err := EmailCheck(newEmail.Email)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "email can't be empty or malformed")
		return
	}
	// update student email
	studentProfile, err := cfg.DB.UpdateStudentEmailByID(r.Context(), database.UpdateStudentEmailByIDParams{Email: email, ID: studentID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("student not found: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "unable to update email")
			return
		}
		// detect the duplicate email error
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			ResponseWithError(w, http.StatusConflict, "email is already in use")
			return
		}
		log.Printf("error updating student email: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to update student email")
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

// let student edit/update their name
func (cfg *ApiConfig) HandlerStudentUpdateName(w http.ResponseWriter, r *http.Request) {
	// get student ID from the context
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// decode edit name request
	var newName UpdateNameRequest
	err := DecodeRequest(r, &newName)
	if err != nil {
		log.Printf("error decoding update name request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to update name")
		return
	}

	// check student's name input
	firstName, err := NameCheck(newName.FirstName)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "name can't be empty")
		return
	}
	lastName, err := NameCheck(newName.LastName)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "name can't be empty")
		return
	}

	// update student name
	studentProfile, err := cfg.DB.UpdateStudentName(r.Context(), database.UpdateStudentNameParams{
		FirstName: firstName,
		LastName:  lastName,
		ID:        studentID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("student not found: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "unable to update name")
			return
		}
		log.Printf("error updating student name: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to update student name")
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

// let student remove a course from their course list
func (cfg *ApiConfig) HandlerStudentRemoveACourse(w http.ResponseWriter, r *http.Request) {
	// check student ID
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// get course id from url
	course_id := r.PathValue("id")
	if course_id == "" {
		ResponseWithError(w, http.StatusBadRequest, "course id can't be empty")
		return
	}

	// remove course from student's course list
	err := cfg.DB.RemoveACourse(r.Context(), database.RemoveACourseParams{StudentID: studentID, CourseID: course_id})
	if err != nil {
		log.Printf("error removing course from student's course list: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to remove course")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Course has been removed",
	})
}
