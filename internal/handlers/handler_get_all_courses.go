package handlers

import (
	"log"
	"net/http"
)

// let students see available courses
func (cfg *ApiConfig) HandlerGetAllCourses(w http.ResponseWriter, r *http.Request) {
	// create a list, for about 10 different courses
	list := make([]Course, 0, 10)

	// get list of courses
	courses, err := cfg.DB.ListAllCourses(r.Context())
	if err != nil {
		log.Printf("error getting list of all courses: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to get list of all courses")
		return
	}

	// write courses to format response
	for _, c := range courses {
		list = append(list, Course{
			ID:         c.ID,
			CourseName: c.CourseName,
			Section:    c.SectionDate,
			Time:       c.StartTime,
		})
	}

	// response
	ResponseWithJSON(w, http.StatusOK, list)
}
