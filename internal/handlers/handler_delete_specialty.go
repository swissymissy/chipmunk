package handlers

import (
	"log"
	"net/http"
	"strconv"
)

// let professor delete a specialty in the list
func (cfg *ApiConfig) HandlerDeleteSpecialty(w http.ResponseWriter, r *http.Request) {
	// get specialty id 
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "invalid specialty id")
		return
	}
	// remove a specialty in the list
	err = cfg.DB.DeleteSpecialty(r.Context(), id)
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