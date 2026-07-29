package handlers

import (
	"log"
	"net/http"
)

// let professor see every student in the database, regardless of course
// enrollment. this is the entry point for managing a student (e.g. resetting a
// forgotten password) after they have been removed from a course roster and are
// therefore no longer reachable from the per-course roster view.
func (cfg *ApiConfig) HandlerListAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := cfg.DB.ListAllStudents(r.Context())
	if err != nil {
		log.Printf("error fetching all students: %s\n", err)
		ResponseWithError(w, http.StatusInternalServerError, "failed to get list of students")
		return
	}

	// map to the API model — ListAllStudents is SELECT *, so each row carries
	// password_hash / verified / must_change_password. Those must NOT be
	// returned to the client, so we copy only the safe fields (same pattern as
	// HandlerRosters).
	list := make([]Student, 0, len(students))
	for _, s := range students {
		list = append(list, Student{
			ID:        s.ID,
			StudentID: s.StudentID,
			Email:     s.Email,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Specialty: s.Specialty.String,
		})
	}

	ResponseWithJSON(w, http.StatusOK, list)
}
