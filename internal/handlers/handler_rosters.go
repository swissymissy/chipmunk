package handlers

import (
	"log"
	"net/http"
)

// show list of all students that have registed in a course
func (cfg *ApiConfig) HandlerRosters(w http.ResponseWriter, r *http.Request) {
	// get course id from url
	courseID := r.PathValue("course_id")
	if courseID == "" {
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// get list of all students registerd for the course
	students, err := cfg.DB.GetAllStudentsByCourse(r.Context(), courseID)
	if err != nil {
		log.Printf("error fetching all students in a course: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "failed to get list of students in course")
		return
	}

	// write to format repsonse
	var list []Student
	for _, s := range students {
		list = append(list, Student{
			ID:        s.ID,
			StudentID: s.StudentID,
			Email:     s.Email,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Specialty: s.Specialty.String,
		})
	}

	// repsonse
	ResponseWithJSON(w, http.StatusOK, list)
}
