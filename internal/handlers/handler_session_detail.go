package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// let professor see the detail of a session to check its correctness
func (cfg *ApiConfig) HandlerSessionDetail(w http.ResponseWriter, r *http.Request) {
	// get session id in url
	sessionIDStr := r.PathValue("id")
	if sessionIDStr == "" {
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error converting session ID string to int64: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// get session detail
	session, err := cfg.DB.GetSessionByID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("attempt fetch non-exist session: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "session not found")
			return
		}
		log.Printf("error fetching session detail: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to fetch session detail")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, Session{
		ID: session.ID,
		CourseID: session.CourseID,
		SessionDate: session.SessionDate,
		Status: session.Status,
		StartedAt: session.StartedAt,
		EndedAt: session.EndedAt.String,
	})
}