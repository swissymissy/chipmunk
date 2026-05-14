package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
)

// let professor remove a course from the list
func (cfg *ApiConfig) HandlerRemoveCourse(w http.ResponseWriter, r *http.Request) {
	// get course id from url
	course_id := r.PathValue("id")
	if course_id == "" {
		ResponseWithError(w, http.StatusBadRequest, "course id can't be empty")
		return
	}

	// remove a course in the list
	err := cfg.DB.DeleteCourse(r.Context(), course_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("error removing course. Course id not found: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "course not found or already deleted")
			return
		}
		log.Printf("error removing a course in list: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to remove course")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Course has been removed.",
	})

}
