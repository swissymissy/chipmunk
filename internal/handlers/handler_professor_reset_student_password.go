package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/auth"
	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor reset student's password and create a temp password
func (cfg *ApiConfig) HandlerProfessorResetStudentPassword(w http.ResponseWriter, r *http.Request) {
	// get student uuid from url
	studentID := r.PathValue("student_id")
	if studentID == "" {
		ResponseWithError(w, http.StatusBadRequest, "student's id can't be empty")
		return
	}

	// decode change password request
	var resetPassword TempPasswordRequest
	err := DecodeRequest(r, &resetPassword)
	if err != nil {
		log.Printf("error decoding reset password request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to reset password")
		return
	}

	// hash temp password
	hash, err := auth.HashPassword(resetPassword.NewPassword)
	if err != nil {
		log.Printf("error hashing new password: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to reset password")
		return
	}

	// update password
	err = cfg.DB.ResetStudentPassword(r.Context(), database.ResetStudentPasswordParams{PasswordHash: ToNullString(hash), ID: studentID})
	if err != nil {
		log.Printf("error updating password: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to reset password. Something went wrong")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Password is updated. Student will have to change their password in next log in",
	})
}
