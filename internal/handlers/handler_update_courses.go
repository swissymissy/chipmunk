package handlers

import (
	"net/http"
)

// let professor update course's name
func (cfg *ApiConfig) HandlerUpdateCourse(w http.ResponseWriter, r *http.Request) {
	// get course id from url
	idStr := r.PathValue("id")

	// decode to get course new name

}
