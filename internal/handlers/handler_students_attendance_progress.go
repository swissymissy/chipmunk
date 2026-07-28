package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/middleware"
)

// let students see their attandance progress during semester
func (cfg *ApiConfig) HandlerStudentsAttendanceProgress(w http.ResponseWriter, r *http.Request) {
	// check student's ID
	studentID, ok := middleware.GetUserID(r.Context())
	if !ok {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// retrieve student's attendance progress from database
	studentProg, err := cfg.DB.GetStudentAttendanceProgress(r.Context(), studentID)
	if err != nil {
		log.Printf("error getting student attendance progress: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to retrieve attendance progress")
		return
	}

	// response
	list := make([]StudentAttendanceProgressResponse, 0)
	for _, s := range studentProg {
		list = append(list, StudentAttendanceProgressResponse{
			CourseID:      s.CourseID,
			CourseName:    s.CourseName,
			Present:       s.Present,
			Late:          s.Late,
			Absent:        s.Absent,
			TotalSessions: s.TotalSessions,
		})
	}

	ResponseWithJSON(w, http.StatusOK, list)
}
