package handlers

import (
	"log"
	"net/http"
)

// reset table specialties
func (cfg *ApiConfig) HandlerResetSpecialty(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// reset table
	err := cfg.DB.ResetSpecialties(r.Context())
	if err != nil {
		log.Printf("error resetting specialties table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset table specialties")
		return
	}

	w.WriteHeader(200)
}
