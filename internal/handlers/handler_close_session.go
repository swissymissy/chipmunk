package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
)

// let professor close an active session after the class is done or check-in is done
func (cfg *ApiConfig) HandlerCloseSession(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req CloseSessionRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding close session request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// close session
	session, err := cfg.DB.CloseSession(r.Context(), req.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ResponseWithError(w, http.StatusNotFound, "session not found")
			return
		}
		log.Printf("error closing session: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to close session")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, CloseSessionResponse{
		SessionID:   session.ID,
		CourseID:    session.CourseID,
		SessionDate: session.SessionDate,
		Status:      session.Status,
		StartedAt:   LocalizeSQLiteTime(session.StartedAt),
		EndedAt:     LocalizeSQLiteTime(session.EndedAt.String),
	})

}
