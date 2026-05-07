package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor change student's status to present
// get request
func (cfg *ApiConfig) HandlerMarkStudentPresent(w http.ResponseWriter, r *http.Request) {
	// decode request
	type request struct {
		StudentID string `json:"student_id"`
		SessionID int64  `json:"session_id"`
	}
	var req request
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding change student's status request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// set student's status to present
	student, err := cfg.DB.UpdateCheckIn(r.Context(), database.UpdateCheckInParams{
		StudentID: req.StudentID,
		SessionID: req.SessionID,
	})
	if err != nil {
		log.Printf("error updating student's status: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to change student's status")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, struct {
		SessionID int64  `json:"session_id"`
		StudentID string `json:"student_id"`
		Status    string `json:"status"`
		CheckInAt string `json:"checkin_at"`
	}{
		SessionID: student.SessionID,
		StudentID: student.StudentID,
		Status:    student.Status,
		CheckInAt: LocalizeSQLiteTime(student.CheckInAt.String),
	})
}
