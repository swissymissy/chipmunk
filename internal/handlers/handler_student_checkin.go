package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/auth"
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

	
}
