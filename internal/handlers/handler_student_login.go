package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/auth"
)

// student login handler
func (cfg *ApiConfig) HandlerStudentLogin(w http.ResponseWriter, r *http.Request) {
	// decode request
	var loginReq StudentLoginRequest
	err := DecodeRequest(r, &loginReq)
	if err != nil {
		log.Printf("error decoding login request: %s\n", err)
		ResponseWithError(w, 400, "invalid request")
		return
	}

	// get student by email
	student, err := cfg.DB.GetStudentByEmail(r.Context(), loginReq.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Login attempt with unknown Email: %s\n", err)
			ResponseWithError(w, 401, "Incorret Email or Password")
			return
		}
		log.Print("error getting student from db: %s\n", err)
		ResponseWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
		return
	}

	// check password
	match, err := auth.CheckPasswordHash(loginReq.Password, student.PasswordHash)
	if err != nil {
		log.Printf("%s\n", err)
		ResponseWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
		return
	}
	if !match {
		ResponseWithError(w, http.StatusUnauthorized, "Incorrect Email or Password")
		return
	}

	
}