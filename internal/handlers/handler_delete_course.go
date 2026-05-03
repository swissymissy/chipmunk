package handlers

import (
	"log"
	"net/http"
)

// reset courses table
// reset this table after resetting enrollments table
func (cfg *ApiConfig) HandlerResetCourses(w http.ResponseWriter, r *http.Request) {
	// check
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := cfg.DB.DeleteAllCourse(r.Context())
	if err != nil {
		log.Printf("error reseting courses table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset table")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Successfully reset all courses.",
	})
}
