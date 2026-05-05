package handlers

import (
	"log"
	"net/http"
)

type specialty struct {
	Name string `json:"specialty_name"`
}

// professor create specialty list
func (cfg *ApiConfig) HandlerCreateSpecialty(w http.ResponseWriter, r *http.Request) {
	// decode req
	var req specialty
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding new specialty request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "failed to add new specialty")
		return
	}

	// check input
	if req.Name == "" {
		ResponseWithError(w, http.StatusBadRequest, "specialty name can't be empty")
		return
	}

	s, err := cfg.DB.CreateSpecialty(r.Context(), req.Name)
	if err != nil {
		log.Printf("error creating new specialty: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to create new specialty")
		return
	}

	ResponseWithJSON(w, http.StatusOK, Specialty{
		ID: s.ID,
		Name: s.Name,
		CreatedAt: s.CreatedAt,
	})
}