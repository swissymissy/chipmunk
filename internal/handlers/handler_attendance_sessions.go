package handlers

import (
	"log"
	"net/http"
	"sort"
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
	Flagged    bool   `json:"flagged"`           // flagged students that has same fingerprint with another student
}

// let professor see all the rosters in the specific session and check their status
// let professor see any device fingerprint conflict
// the frontend will poll for this request every 5 seconds
func (cfg *ApiConfig) HandlerAttendanceBySession(w http.ResponseWriter, r *http.Request) {
	// get session id from URL
	sessionIDStr := r.PathValue("session_id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error converting string to int: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	// get list of all students in the current session
	rosters, err := cfg.DB.GetRecordBySession(r.Context(), sessionID)
	if err != nil {
		log.Printf("error fetching students in the courses: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to get rosters in session")
		return
	}

	// flagged rows, contains list of students
	// that have fingerprint matches at least 1 other student in current session
	flaggedrows, err := cfg.DB.GetFlaggedFingerprints(r.Context(), sessionID)
	if err != nil {
		log.Printf("error fetching flagged fingerprint rows: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to fetch flags")
		return
	}

	// group students up by fingerprints
	// create a set of flagged students for frontend to show warning badge
	groupByFinger := make(map[string][]FlaggedStudent) // DS looks like: [fingerprint key]:[ {}FlaggedStudent, {}FlaggedStudent, {}FlaggedStudent, etc.]
	flagged := make(map[string]bool)
	for _, f := range flaggedrows {
		fingerprint := f.DeviceFingerprint.String
		groupByFinger[fingerprint] = append(groupByFinger[fingerprint], FlaggedStudent{
			StudentID: f.StudentID,
			SchoolID:  f.SchoolID,
			FirstName: f.FirstName,
			LastName:  f.LastName,
			CheckInAt: LocalizeSQLiteTime(f.CheckInAt.String),
		})
		flagged[f.StudentID] = true
	}

	// create an array that has all groups of flagged students
	flagGroup := make([]FlagGroups, 0, len(groupByFinger))
	for key, val := range groupByFinger {
		flagGroup = append(flagGroup, FlagGroups{
			Fingerprint: key,
			Students:    val,
		})
	}

	// sort the list to make stable look for UI
	// sort by fingerprint
	sort.Slice(flagGroup, func(i, j int) bool {
		return flagGroup[i].Fingerprint < flagGroup[j].Fingerprint
	})

	// build roster list
	// mark flagged students
	rosterList := make([]RosterRep, 0, len(rosters))
	for _, s := range rosters {
		rosterList = append(rosterList, RosterRep{
			SessionID:  s.SessionID,
			StudentID:  s.StudentID,
			Status:     s.Status,
			CheckInAt:  LocalizeSQLiteTime(s.CheckInAt.String),
			FirstName:  s.FirstName,
			LastName:   s.LastName,
			StudentID2: s.StudentID_2,
			Flagged:    flagged[s.StudentID], // return true or false
		})
	}

	// response
	ResponseWithJSON(w, http.StatusOK, AttendanceBySessionResponse{
		Roster:     rosterList,
		FlagGroups: flagGroup,
	})
}
