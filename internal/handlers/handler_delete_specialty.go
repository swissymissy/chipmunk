package handlers

import (
	"log"
	"net/http"
)

type req struct {
	ID int64 `json:"id"`
	Name string `json:"specialty_name"`
}

// let professor delete a specialty in the list
func (cfg *ApiConfig) HandlerDeleteSpecialty(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req req
	err := DecodeRequest(r, &req)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "failed to remove specialty")
		return
	}

	// remove a specialty in the list
	err = cfg.DB.DeleteSpecialty(r.Context(), req.ID)
	if err != nil {
		log.Printf("error removing a specialty: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to remove a specialty")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct{
		Msg string `json:"msg"`
	}{
		Msg: "Specialty has been removed.",
	})
}