package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor edit the name of a specialty
// it will also check the unique of the specialty name
func (cfg *ApiConfig) HandlerUpdateSpecialty(w http.ResponseWriter, r *http.Request) {
	// get specialty id
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("error converting string to int64: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid specialty id")
		return
	}

	// decode to get new specialty name
	var specialtyName UpdateSpecialtyReq
	err = DecodeRequest(r, &specialtyName)
	if err != nil {
		log.Printf("error decoding request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to update specialty name")
		return
	}

	// check input
	newName, err := NameCheck(specialtyName.Name)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "name can't be empty")
		return
	}
	// check specialty length
	if !MaxLenOK(newName, 128) {
		ResponseWithError(w, http.StatusBadRequest, "name is too long")
		return
	}

	// update name in database
	newSpecialty, err := cfg.DB.UpdateSpecialty(r.Context(), database.UpdateSpecialtyParams{ID: id, Name: newName})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("error updating specialty name: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "specialty not found")
			return
		}

		// detect duplicate specialty
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			ResponseWithError(w, http.StatusConflict, "specialty already in use")
			return
		}

		log.Printf("error updating specialty name: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to update specialty")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, struct {
		Msg  string `json:"msg"`
		Name string `json:"name"`
	}{
		Msg:  "Specialty has been updated",
		Name: newSpecialty.Name,
	})
}
