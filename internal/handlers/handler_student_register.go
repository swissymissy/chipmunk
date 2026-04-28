package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/chipmunk/internal/auth"
	"github.com/swissymissy/chipmunk/internal/database"
)

// let students register the first time
func (cfg *ApiConfig) HandlerStudentRegister(w http.ResponseWriter, r *http.Request) {
	// decode register request
	var req StudentRegisterRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding register request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to create student")
		return
	}

	// check input
	if req.StudentID == "" || req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Specialty == "" {
		ResponseWithError(w, http.StatusBadRequest, "please fill up required information")
		return
	}

	// hash password
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("error hashing student's password: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// create new uuid
	studentUID := uuid.New().String()

	// create new student in database
	student, err := cfg.DB.CreateStudent(r.Context(), database.CreateStudentParams{
		ID:           studentUID,
		StudentID:    req.StudentID,
		Email:        req.Email,
		PasswordHash: ToNullString(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Specialty:    ToNullString(req.Specialty),
	})
	if err != nil {
		log.Printf("error creating new student: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to create account")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, StudentRegisterResponse{
		StudentID: student.StudentID,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
		Specialty: student.Specialty.String,
	})
}
