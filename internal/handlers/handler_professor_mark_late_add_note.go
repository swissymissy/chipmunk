package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/database"
)

// let professor mark a student to be late in past session
// let professor add note to an absent student
func (cfg *ApiConfig) HandlerUpdateAttendanceRecord(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req UpdateAttendanceRecordRequest
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding update attendance request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "unable to update record")
		return
	}

	// check input
	if req.Status != "late" && req.Status != "present" && req.Status != "absent" {
		log.Printf("invalid status: %s\n", req.Status)
		ResponseWithError(w, http.StatusBadRequest, "invalid status")
		return
	}
	if !MaxLenOK(req.Note, 500) {
		log.Println("note is too long")
		ResponseWithError(w, http.StatusBadRequest, "note is too long")
		return
	}

	// update record in database
	record, err := cfg.DB.UpdateAttendanceRecord(r.Context(), database.UpdateAttendanceRecordParams{
		Status:    req.Status,
		Note:      ToNullString(req.Note),
		SessionID: req.SessionID,
		StudentID: req.StudentID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("record is not found: %s\n", err)
			ResponseWithError(w, http.StatusNotFound, "unable to update record. Record not found.")
			return
		}
		log.Printf("error updating record for student: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to update student attendance record")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, AttendanceRecordResponse{
		Msg:       "Record has been updated successfully.",
		SessionID: record.SessionID,
		StudentID: record.StudentID,
		Status:    record.Status,
		Note:      record.Note.String,
	})

}
