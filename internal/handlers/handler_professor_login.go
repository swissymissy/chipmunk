package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/auth"
)

type ProfessorLoginReq struct {
	Password string `json:"password"`
}

type ProfessorLoginRes struct {
	Token string `json:"token"`
}

// let professor login to dashboard
func (cfg *ApiConfig) HandlerProfessorLogin(w http.ResponseWriter, r *http.Request) {
	var req ProfessorLoginReq
	err := DecodeRequest(r, &req)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if cfg.ProfessorPasswordHash == "" {
		log.Printf("PROFESSOR_PASSWORD_HASH not set")
		ResponseWithError(w, http.StatusInternalServerError, "server not configured")
		return
	}

	match, err := auth.CheckPasswordHash(req.Password, cfg.ProfessorPasswordHash)
	if err != nil {
		log.Printf("password check error: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
	} 
	if !match {
		ResponseWithError(w, http.StatusUnauthorized, "incorrect password")
		return
	}

	token, err := auth.MakeProfessorJWT(cfg.JWT)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	log.Println("Professor has logged in")
	ResponseWithJSON(w, http.StatusOK, ProfessorLoginRes{
		Token: token,
	})
}
