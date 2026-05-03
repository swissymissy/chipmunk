package handlers

import (
	"log"
	"net/http"
)

// reset attendance_sessions table
// reset this table after resetting attendance_record table
func (cfg *ApiConfig) HandlerResetSessions(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := cfg.DB.ResetAttendanceSession(r.Context()); err != nil {
		log.Printf("error resetting attendance_sessions table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset sessions table")
		return
	}
	w.WriteHeader(200)
}
