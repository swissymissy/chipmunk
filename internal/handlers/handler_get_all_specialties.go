package handlers

import (
	"log"
	"net/http"
)

// list all specialties for students to pick during registration
func (cfg *ApiConfig) HandlerGetAllSpecialties(w http.ResponseWriter, r *http.Request) {
	// make a list for about 20 different majors
	list := make([]Specialty, 0, 20)

	// get list of call specialties
	specialties, err := cfg.DB.ListAllSpecialties(r.Context())
	if err != nil {
		log.Printf("error fetching list of all specialties: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to fetch all specialties")
		return
	}

	// write in format response
	for _, s := range specialties {
		list = append(list, Specialty{
			ID:   s.ID,
			Name: s.Name,
		})
	}

	// response
	ResponseWithJSON(w, http.StatusOK, list)
}
