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

	// get qr token from url
	qrToken := r.URL.Query().Get("t")

	// extract sessionID and timestamp from token
	sessionID, timestamp, err := auth.ExtractSessionIDTimestamp(qrToken)
	if err != nil {
		log.Printf("error extracting session id and timestamp from qr token: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid token")
		return
	}

	// get sessionID secret and check if session is active
	session, err := cfg.DB.GetSessionByID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Panicf("attempt to get non-exist session: %s\n", err)
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

	// verify qrToken
	sessionID, err := auth.ValidateQRToken(qrToken, )

	ResponseWithError(w, http.StatusNotImplemented, "not implemented yet")
}
