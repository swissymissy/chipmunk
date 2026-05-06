package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/swissymissy/chipmunk/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

// check student's token and validate token
func AuthRequired(next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for token bearer header
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("No bearer token: %s\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate token
		studentID, err := auth.ValidateJWT(token, jwtSecret)
		if err != nil {
			log.Printf("invalid token: %s\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// store studentID in contex so handlers can access it
		ctx := context.WithValue(r.Context(), UserIDKey, studentID)

		// valid token, let request go through
		log.Printf("student with ID: %s , has logged in\n", studentID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


// protect professor's endpoints
func RequireProfessor(next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		// check for token bearer header
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("No bearer token: %s\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		subject, err := auth.ValidateJWT(token, jwtSecret)
		if err != nil || subject != "professor" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return 
		}
		next.ServeHTTP(w, r)
	})
}