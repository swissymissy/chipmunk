package handlers

import (
	"log"
	"net/http"
	"strconv"
)

type RosterRep struct {
	SessionID  int64  `json:"session_id"`
	StudentID  string `json:"student_id"`
	Status     string `json:"status"`
	CheckInAt  string `json:"checkin_at"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	StudentID2 string `json:"student_school_id"` // student's school id
}

// let professor see all the rosters in the specific session and check their status
// get req
func (cfg *ApiConfig) HandlerAttendanceBySession(w http.ResponseWriter, r *http.Request) {
	// get session id from URL
	sessionIDStr := r.PathValue("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error converting string to int: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// get list of students in the current session
	rosters, err := cfg.DB.GetRecordBySession(r.Context(), sessionID)
	if err != nil {
		log.Printf("error fetching students in the courses: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to get rosters in session")
		return
	}

	// write to format response
	var list []RosterRep
	for _, s := range rosters {
		list = append(list, RosterRep{
			SessionID:  s.SessionID,
			StudentID:  s.StudentID,
			Status:     s.Status,
			CheckInAt:  s.CheckInAt.String,
			FirstName:  s.FirstName,
			LastName:   s.LastName,
			StudentID2: s.StudentID_2,
		})
	}

	// response
	ResponseWithJSON(w, http.StatusOK, list)
}
