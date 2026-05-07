package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/middleware"
)

// list all courses that a student has enrolled in
func (cfg *ApiConfig) HandlerStudentEnrollments(w http.ResponseWriter, r *http.Request) {
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	courseList, err := cfg.DB.GetEnrollmentsByStudent(r.Context(), studentID)
	if err != nil {
		log.Printf("error fetching enrollments for student: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to load enrollments")
		return
	}

	list := make([]Course, 0, 10)
	for _, c := range courseList {
		list = append(list, Course{
			ID:         c.ID,
			CourseName: c.CourseName,
			Section:    c.SectionDate,
			Time:       c.StartTime,
		})
	}

	ResponseWithJSON(w, http.StatusOK, list)
}
