package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/swissymissy/chipmunk/internal/database"
)

// handle the request of adding students to past sessions
func (cfg *ApiConfig) HandlerAddStudentToPastSessionRecord(w http.ResponseWriter, r *http.Request) {
	sessionIDStr := r.PathValue("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error converting session id to int: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	// decode request
	var student AddStudentToPastSessionRecordReq
	err = DecodeRequest(r, &student)
	if err != nil {
		log.Printf("error decoding student id: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	// check body input
	if student.StudentID == "" || !MaxLenOK(student.StudentID, 100) {
		log.Print("student id can't be empty or too long \n")
		ResponseWithError(w, http.StatusBadRequest, "invalid student id")
		return
	}

	// call query to db
	record, err := cfg.DB.AddStudentToSession(r.Context(), database.AddStudentToSessionParams{ID: sessionID, StudentID: student.StudentID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Print("error adding student to session. Student isn't enrolled in this course or already in this session")
			ResponseWithError(w, http.StatusBadRequest, "unable to add student to session. Student isn't enrolled in this course or already in this session")
			return
		}
		log.Printf("error adding student to this session: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "unable to add student to this session.")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusCreated, AddStudentToPastSessionRecordRes{
		Msg:       "Student has been added to session successfully",
		SessionID: record.SessionID,
		StudentID: record.StudentID,
		Status:    record.Status,
		Note:      record.Note.String,
	})
}
