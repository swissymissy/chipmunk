package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/swissymissy/chipmunk/internal/auth"
)

// this function poll for fresh qr token
func (cfg *ApiConfig) HandlerGetQRToken(w http.ResponseWriter, r *http.Request) {
	// get session id from url
	sessionIDStr := r.PathValue("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error parsing session ID: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// fetch session ID to get secret key
	session, err := cfg.DB.GetSessionByID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("attempt fetching non-exist session: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "session not found")
			return
		}
		log.Printf("error fetching session by id: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// check if session is active
	if session.Status != "active" {
		ResponseWithError(w, http.StatusBadRequest, "session is no longer active")
		return
	}

	// generate fresh token
	qrToken, err := auth.GenerateQRToken(sessionID, session.SecretKey)
	if err != nil {
		log.Printf("error generating new qr token: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to generate new qr token")
		return
	}

	// crafting checkin url with token
	checkinURL := fmt.Sprintf("%s/checkin.html?t=%s", cfg.BaseURL, qrToken)

	// response
	ResponseWithJSON(w, http.StatusOK, QRTokenResponse{
		Token:      qrToken,
		CheckInURL: checkinURL,
	})
}
