package handlers

import (
	"log"
	"net/http"
)

// reset students table
// rest this table in last order
func (cfg *ApiConfig) HandlerResetStudents(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := cfg.DB.ResetStudents(r.Context()); err != nil {
		log.Printf("error resetting students table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset students table")
		return
	}

	w.WriteHeader(http.StatusOK)
}