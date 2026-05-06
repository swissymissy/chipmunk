package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/swissymissy/chipmunk/internal/database"
	"github.com/swissymissy/chipmunk/internal/handlers"
	"github.com/swissymissy/chipmunk/internal/middleware"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	baseURL := os.Getenv("BASE_URL")

	dbURL := os.Getenv("DB_URL")
	// open connection to database
	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %s\n", err)
		return
	}
	// set up PRAGMA
	pragmas := []string{
		"PRAGMA journal_mode=WAL",  // write-ahead logging for better concurrency
		"PRAGMA busy_timeout=5000", // sleep for 5s
		"PRAGMA foreign_keys=ON",   // enable foreign keys
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Fatalf("failed to set %s: %v", p, err)
		}
	}
	// query
	dbQuery := database.New(db)
	log.Print("Database connected")

	platform := os.Getenv("PLATFORM")
	jwt := os.Getenv("JWT_SECRET")

	// server config
	cfg := &handlers.ApiConfig{
		Port:     port,
		DB:       dbQuery,
		Platform: platform,
		JWT:      jwt,
		BaseURL:  baseURL,
	}

	// server mux
	mux := http.NewServeMux()

	// new http server
	address := fmt.Sprintf(":%s", port)
	chipmunkServer := http.Server{
		Addr:    address,
		Handler: mux,
	}

	homepage := http.FileServer(http.Dir("./cmd/frontend"))
	mux.Handle("/", homepage)

	// TODO: register handlers
	mux.HandleFunc("GET /api/health", handlers.HandlerHealthCheck)

	// professor - local only
	mux.HandleFunc("POST /api/courses", middleware.LocalOnly(cfg.HandleCreateCourse))                              // professor create new course
	mux.HandleFunc("POST /api/sessions/start", middleware.LocalOnly(cfg.HandlerStartSession))                      // start a new session
	mux.HandleFunc("PUT /api/sessions/close", middleware.LocalOnly(cfg.HandlerCloseSession))                       // close an active session
	mux.HandleFunc("PUT /api/sessions/reopen", middleware.LocalOnly(cfg.HandlerReopenSession))                     // reopen a closed session
	mux.HandleFunc("GET /api/sessions/{id}", middleware.LocalOnly(cfg.HandlerSessionDetail))                       // get session details just to check
	mux.HandleFunc("GET /api/roster/{course_id}", middleware.LocalOnly(cfg.HandlerRosters))                        // view students enrolled in a course
	mux.HandleFunc("GET /api/attendance/{session_id}", middleware.LocalOnly(cfg.HandlerAttendanceBySession))       // view who is present/absent in a specific session
	mux.HandleFunc("PUT /api/attendance/override", middleware.LocalOnly(cfg.HandlerMarkStudentPresent))            // manually mark a student present
	mux.HandleFunc("GET /api/sessions/{id}/qr", middleware.LocalOnly(cfg.HandlerGetQRToken))                       // endpoint for professor to get fresh qr token
	mux.HandleFunc("GET /api/export/semester/{course_id}", middleware.LocalOnly(cfg.HandlerExportSemesterRecords)) // export semester attendance records to excel file
	mux.HandleFunc("GET /api/export/daily/{date}", middleware.LocalOnly(cfg.HandlerExportDailyRecord))             // export daily attendance records to excel file
	mux.HandleFunc("POST /api/specialties", middleware.LocalOnly(cfg.HandlerCreateSpecialty))                      // professor create new specialty
	mux.HandleFunc("DELETE /api/specialties/{id}", middleware.LocalOnly(cfg.HandlerDeleteSpecialty))               // professor delete a specialty in the list
	mux.HandleFunc("GET /api/sessions/active", middleware.LocalOnly(cfg.HandlerListActiveSession))                 // get all active sessions to let professor close in case forget to to close session

	// students - public
	mux.HandleFunc("GET /api/courses", cfg.HandlerGetAllCourses)          // list courses to let students pick
	mux.HandleFunc("POST /api/auth/login", cfg.HandlerStudentLogin)       // student login
	mux.HandleFunc("POST /api/auth/register", cfg.HandlerStudentRegister) // student register new account
	mux.HandleFunc("GET /api/specialties", cfg.HandlerGetAllSpecialties)  // let student see list of all specialties

	// students - auth required
	mux.HandleFunc("POST /api/enrollment", middleware.AuthRequired(cfg.HandlerEnrollment, cfg.JWT))             // student enroll in a course
	mux.HandleFunc("POST /api/attendance/checkin", middleware.AuthRequired(cfg.HandlerStudentCheckIn, cfg.JWT)) // students check in

	// reset - only dev or prof
	mux.HandleFunc("DELETE /api/reset/students", middleware.LocalOnly(cfg.HandlerResetStudents))       // reset students table
	mux.HandleFunc("DELETE /api/reset/courses", middleware.LocalOnly(cfg.HandlerResetCourses))         // reset courses table
	mux.HandleFunc("DELETE /api/reset/enrollments", middleware.LocalOnly(cfg.HandlerResetEnrollments)) // reset enrollments table
	mux.HandleFunc("DELETE /api/reset/sessions", middleware.LocalOnly(cfg.HandlerResetSessions))       // reset attendance sessions table
	mux.HandleFunc("DELETE /api/reset/records", middleware.LocalOnly(cfg.HandlerResetRecords))         // reset attendance records table
	mux.HandleFunc("DELETE /api/reset/all", middleware.LocalOnly(cfg.HandlerResetAll))                 // reset all tables in correct order
	mux.HandleFunc("DELETE /api/reset/specialties", middleware.LocalOnly(cfg.HandlerResetSpecialty))   // reset table specialties

	// run server in background
	go func() {
		fmt.Printf("Serving on: %s:%s/\n", baseURL, port)
		if err := chipmunkServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	// graceful shutdown
	// block until os sends SIGTERM or SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := chipmunkServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error. Forced shutdown: %s\n", err)
	}
	log.Println("Graceful shutdown complete")
}
