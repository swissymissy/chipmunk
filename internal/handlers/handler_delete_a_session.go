package handlers

import (
	"log"
	"net/http"
	"strconv"
)

// let professor delete a session from session list
func (cfg *ApiConfig) HandlerDeleteASession(w http.ResponseWriter, r *http.Request) {
	// get session id from url path
	sessionIDStr := r.PathValue("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error converting session id string to int64: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	// remove session from database
	err = cfg.DB.DeleteSession(r.Context(), sessionID)
	if err != nil {
		log.Printf("error deleting session from database: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to delete session")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Session has been removed.",
	})
}
