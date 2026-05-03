package handlers

import (
	"log"
	"net/http"
)

// reset enrollments table
// reset this table after restting attendance_sessions table
func (cfg *ApiConfig) HandlerResetEnrollments(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := cfg.DB.ResetEnrollment(r.Context()); err != nil {
		log.Printf("error resetting enrollments table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset enrollments table")
		return
	}
	w.WriteHeader(http.StatusOK)
}