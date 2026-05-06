package handlers

import (
	"log"
	"net/http"
)

// list all current active sessions
// in case professor forgot to close a session
func (cfg *ApiConfig) HandlerListActiveSession(w http.ResponseWriter, r *http.Request) {

	session, err := cfg.DB.ListActiveSessions(r.Context())
	if err != nil {
		log.Printf("error fetching active sessions: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	list := make([]ActiveSessionsView, 0, len(session))
	for _, s := range session {
		list = append(list, ActiveSessionsView{
			SessionID:   s.ID,
			CourseID:    s.CourseID,
			SessionDate: s.SessionDate,
			Status:      s.Status,
			StartedAt:   s.StartedAt,
		})
	}
	ResponseWithJSON(w, http.StatusOK, list)
}
