package handlers

import (
	"log"
	"net/http"
	"strconv"
)

// let professor see the detail of a session to check its correctness
func (cfg *ApiConfig) HandlerSessionDetail(w http.ResponseWriter, r *http.Request) {
	// get session id in url
	sessionIDStr := r.PathValue("id")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		log.Printf("error converting session ID string to int64: %s\n", err)
		ResponseWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	
}