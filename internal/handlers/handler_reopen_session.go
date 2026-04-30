package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
)

type ReopenSessionRequest struct {
	SessionID int64 `json:"session_id"`
}

type ReopenSessionResponse struct {
	SessionID   int64  `json:"session_id"`
	CourseID    string `json:"course_id"`
	SessionDate string `json:"session_date"`
	Status      string `json:"status"`
	StartedAt   string `json:"started_at"`
}

// let professor reopen a closed session
func (cfg *ApiConfig) HandlerReopenSession(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req ReopenSessionRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding reopen session request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// reopen session
	session, err := cfg.DB.ReOpenSession(r.Context(), req.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("attempt to reopen a non-exist session: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "session does not exist")
			return
		}
		log.Printf("error reopening session: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reopen session")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, ReopenSessionResponse{
		SessionID:   session.ID,
		CourseID:    session.CourseID,
		SessionDate: session.SessionDate,
		Status:      session.Status,
		StartedAt:   session.StartedAt,
	})
}
