package handlers

import (
	"log"
	"net/http"
)

// let professor delete a student account from database
func (cfg *ApiConfig) HandlerDeleteAStudentAccount(w http.ResponseWriter, r *http.Request) {
	// get student uuid from url
	stdID := r.PathValue("id")

	// remove student account from database
	err := cfg.DB.DeleteStudent(r.Context(), stdID)
	if err != nil {
		log.Printf("error deleting a student account from database: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to delete student account")
		return
	}

	// response
	ResponseWithJSON(w, http.StatusOK, struct {
		Msg string `json:"msg"`
	}{Msg: "Student account has been removed."})
}
