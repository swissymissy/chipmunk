package handlers

import (
	"log"
	"net/http"
)

type RosterReq struct {
	SessionID int64 `json:"session_id"`
}
type RosterRep struct {
	Status    string `json:"status"`
	CheckInAt string `json:"checkin_at"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	StudentID string `json:"student_id"` // student's school id
}

// let professor see all the rosters in a course
func (cfg *ApiConfig) HandlerRosters(w http.ResponseWriter, r *http.Request) {
	// decode request
	var req RosterReq
	err := DecodeRequest(r, &req)
	if err != nil {
		log.Printf("error decoding rosters request: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// get list of students in the course
	rosters, err := cfg.DB.GetRecordBySession(r.Context(), req.SessionID)
	if err != nil {
		log.Printf("error fetching students in the courses: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to get rosters in course")
		return
	}

	// write to format response
	var list []RosterRep
	for _, s := range rosters {
		list = append(list, RosterRep{
			Status:    s.Status,
			CheckInAt: s.CheckInAt.String,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			StudentID: s.StudentID_2,
		})
	}

	// response
	ResponseWithJSON(w, http.StatusOK, list)
}
