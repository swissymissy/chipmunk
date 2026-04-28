package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor create new course
func (cfg *ApiConfig) HandleCreateCourse(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req NewCourseRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding new course request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to create new course")
		return
	}

	// check input
	if req.Name == "" || req.Section == "" || req.Time == "" {
		ResponseWithError(w, http.StatusBadRequest, "please make sure to fill up all fields")
		return
	}

	courseID := uuid.New().String()

	// insert new course to database
	course, err := cfg.DB.CreateCourse(r.Context(), database.CreateCourseParams{
		ID:          courseID,
		CourseName:  req.Name,
		SectionDate: req.Section,
		StartTime:   req.Time,
	})
	if err != nil {
		log.Printf("error inserting new course to database: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to create course")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, NewCourseResponse{
		ID:      course.ID,
		Name:    course.CourseName,
		Section: course.SectionDate,
		Time:    course.StartTime,
	})
}
