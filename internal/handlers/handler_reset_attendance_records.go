package handlers

import (
	"log"
	"net/http"
)

// reset table attendance_records
// rest this table first
func (cfg *ApiConfig) HandlerResetRecords(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := cfg.DB.ResetAttendanceRecords(r.Context()); err != nil {
		log.Printf("error resetting attendance records table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset attendance records table")
		return
	}

	w.WriteHeader(200)
}