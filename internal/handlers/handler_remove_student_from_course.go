// let professor remove a student from a course
package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
)

func (cfg *ApiConfig) HandlerRemoveStudentFromCourse(w http.ResponseWriter, r *http.Request) {
	// get course id from url
	course_id := r.PathValue("course_id")
	if course_id == "" {
		ResponseWithError(w, http.StatusBadRequest, "course id can't be empty")
		return
	}
	// get student uuid from url
	studentID := r.PathValue("student_id")
	if studentID == "" {
		ResponseWithError(w, http.StatusBadRequest, "student's id can't be empty")
		return
	}

	// remove student from course
	err := cfg.DB.RemoveACourse(r.Context(), database.RemoveACourseParams{StudentID: studentID, CourseID: course_id})
	if err != nil {
		log.Printf("error removing student from course: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to remove student from course")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "Student has been removed from course",
	})
}
