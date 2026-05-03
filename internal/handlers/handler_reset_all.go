package handlers

import (
	"log"
	"net/http"
)

// reset all tables in reverse dependency order
// at the end of semester, professor can reest all data to prepare for new semester
func (cfg *ApiConfig) HandlerResetAll(w http.ResponseWriter, r *http.Request) {
	// check
	if cfg.Platform != "dev" && cfg.Platform != "prof" {
		ResponseWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 1. reset attendance_records table
	err := cfg.DB.ResetAttendanceRecords(r.Context())
	if err != nil {
		log.Printf("error reseting attendance_records: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset records tables")
		return
	}

	// 2. reset attedance_sessions table
	err = cfg.DB.ResetAttendanceSession(r.Context())
	if err != nil {
		log.Printf("error resetting attendance_sessions: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset sessions table")
		return
	}

	// 3. reset enrollments table
	if err := cfg.DB.ResetEnrollment(r.Context()); err != nil {
		log.Printf("error resetting enrollments table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset enrollments table")
		return
	}

	// 4. reset courses table
	if err := cfg.DB.DeleteAllCourse(r.Context()); err != nil {
		log.Printf("error resetting courses table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset courses table")
		return
	}

	// 5. reset students table
	if err := cfg.DB.ResetStudents(r.Context()); err != nil {
		log.Printf("error resetting students table: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to reset students table")
		return
	}

	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"message"`
	}{
		Msg: "Successfully reset all tables",
	})
}
