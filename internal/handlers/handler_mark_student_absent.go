package handlers

import (
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor flip student's status back to 'absent'
// can be used after reviewing a flagged
func (cfg *ApiConfig) HandlerMarkStudentAbsent(w http.ResponseWriter, r *http.Request) {
	// decode request
	type request struct {
		StudentID string `json:"student_id"`
		SessionID int64  `json:"session_id"`
	}
	var req request
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// set student status to 'absent'
	student, err := cfg.DB.RevertCheckin(r.Context(), database.RevertCheckinParams{
		StudentID: req.StudentID,
		SessionID: req.SessionID,
	})
	if err != nil {
		log.Printf("error updating student status to 'absent': %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to update student status")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, struct {
		SessionID int64  `json:"session_id"`
		StudentID string `json:"student_id"`
		Status    string `json:"status"`
	}{
		SessionID: student.SessionID,
		StudentID: student.StudentID,
		Status:    student.Status,
	})
}
