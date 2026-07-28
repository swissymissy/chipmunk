package handlers

import (
	"log"
	"net/http"
)

// let professor see all past sessions by a course
func (cfg *ApiConfig) HandlerListAllSessions(w http.ResponseWriter, r *http.Request) {
	// get course id from url
	courseID := r.PathValue("course_id")

	// retrieve all past sessions of the course
	listSessions, err := cfg.DB.ListAllSessionsByCourse(r.Context(), courseID)
	if err != nil {
		log.Printf("error getting all sessions of the course: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to get all sessions for this course")
		return
	}

	// response
	list := make([]AllSessionByCourseResponse, 0)
	for _, sess := range listSessions {
		list = append(list, AllSessionByCourseResponse{
			ID:          sess.ID,
			SessionDate: sess.SessionDate,
			Status:      sess.Status,
			StartedAt:   sess.StartedAt,
			EndedAt:     sess.EndedAt.String,
		})
	}

	ResponseWithJSON(w, http.StatusOK, list)
}
