package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor update course's information: name, section date, time
func (cfg *ApiConfig) HandlerUpdateCourse(w http.ResponseWriter, r *http.Request) {
	// get course id from url
	idStr := r.PathValue("id")

	// decode to get course new name
	var course UpdateCourseReq
	err := DecodeRequest(r, &course)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// check input
	name, err := NameCheck(course.CourseName)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "name can't be empty")
		return
	}
	section, err := NameCheck(course.Section)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "section date can't be empty")
		return
	}
	time, err := NameCheck(course.Time)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "time can't be empty")
		return
	}
	if !MaxLenOK(name, 128) || !MaxLenOK(section, 128) || !MaxLenOK(time, 128) {
		ResponseWithError(w, http.StatusBadRequest, "name or section or time is too long")
		return
	}

	// update course infor in database
	updateCourse, err := cfg.DB.UpdateCourse(r.Context(), database.UpdateCourseParams{
		ID:          idStr,
		CourseName:  name,
		SectionDate: section,
		StartTime:   time,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("error updating course. course not found: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "course not found")
			return
		}
		log.Printf("error updating course: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to update course information")
		return
	}
	// response
	ResponseWithJSON(w, http.StatusOK, struct {
		Msg        string `json:"msg"`
		CourseName string `json:"course_name"`
		Section    string `json:"section_date"`
		Time       string `json:"start_time"`
	}{
		Msg:        "Course has been updated",
		CourseName: updateCourse.CourseName,
		Section:    updateCourse.SectionDate,
		Time:       updateCourse.StartTime,
	})

}
