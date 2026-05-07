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
	professorHash := os.Getenv("PROFESSOR_PASSWORD_HASH")

	// server config
	cfg := &handlers.ApiConfig{
		Port:                  port,
		DB:                    dbQuery,
		Platform:              platform,
		JWT:                   jwt,
		BaseURL:               baseURL,
		ProfessorPasswordHash: professorHash,
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

	// register handlers
	mux.HandleFunc("GET /api/health", handlers.HandlerHealthCheck)
	mux.HandleFunc("POST /api/auth/professor/login", cfg.HandlerProfessorLogin)

	// professor only
	mux.HandleFunc("POST /api/courses", middleware.RequireProfessor(cfg.HandleCreateCourse, cfg.JWT))                              // professor create new course
	mux.HandleFunc("POST /api/sessions/start", middleware.RequireProfessor(cfg.HandlerStartSession, cfg.JWT))                      // start a new session
	mux.HandleFunc("PUT /api/sessions/close", middleware.RequireProfessor(cfg.HandlerCloseSession, cfg.JWT))                       // close an active session
	mux.HandleFunc("PUT /api/sessions/reopen", middleware.RequireProfessor(cfg.HandlerReopenSession, cfg.JWT))                     // reopen a closed session
	mux.HandleFunc("GET /api/sessions/{id}", middleware.RequireProfessor(cfg.HandlerSessionDetail, cfg.JWT))                       // get session details just to check
	mux.HandleFunc("GET /api/roster/{course_id}", middleware.RequireProfessor(cfg.HandlerRosters, cfg.JWT))                        // view students enrolled in a course
	mux.HandleFunc("GET /api/attendance/{session_id}", middleware.RequireProfessor(cfg.HandlerAttendanceBySession, cfg.JWT))       // view who is present/absent in a specific session
	mux.HandleFunc("PUT /api/attendance/override", middleware.RequireProfessor(cfg.HandlerMarkStudentPresent, cfg.JWT))            // manually mark a student present
	mux.HandleFunc("GET /api/sessions/{id}/qr", middleware.RequireProfessor(cfg.HandlerGetQRToken, cfg.JWT))                       // endpoint for professor to get fresh qr token
	mux.HandleFunc("GET /api/export/semester/{course_id}", middleware.RequireProfessor(cfg.HandlerExportSemesterRecords, cfg.JWT)) // export semester attendance records to excel file
	mux.HandleFunc("GET /api/export/daily/{date}", middleware.RequireProfessor(cfg.HandlerExportDailyRecord, cfg.JWT))             // export daily attendance records to excel file
	mux.HandleFunc("POST /api/specialties", middleware.RequireProfessor(cfg.HandlerCreateSpecialty, cfg.JWT))                      // professor create new specialty
	mux.HandleFunc("DELETE /api/specialties/{id}", middleware.RequireProfessor(cfg.HandlerDeleteSpecialty, cfg.JWT))               // professor delete a specialty in the list
	mux.HandleFunc("GET /api/sessions/active", middleware.RequireProfessor(cfg.HandlerListActiveSession, cfg.JWT))                 // get all active sessions to let professor close in case forget to to close session

	// students - public
	mux.HandleFunc("GET /api/courses", cfg.HandlerGetAllCourses)          // list courses to let students pick
	mux.HandleFunc("POST /api/auth/login", cfg.HandlerStudentLogin)       // student login
	mux.HandleFunc("POST /api/auth/register", cfg.HandlerStudentRegister) // student register new account
	mux.HandleFunc("GET /api/specialties", cfg.HandlerGetAllSpecialties)  // let student see list of all specialties

	// students - auth required
	mux.HandleFunc("POST /api/enrollment", middleware.AuthRequired(cfg.HandlerEnrollment, cfg.JWT))             // student enroll in a course
	mux.HandleFunc("POST /api/attendance/checkin", middleware.AuthRequired(cfg.HandlerStudentCheckIn, cfg.JWT)) // students check in
	mux.HandleFunc("GET /api/enrollments", middleware.AuthRequired(cfg.HandlerStudentEnrollments, cfg.JWT))     // show list of all courses student has enrolled in

	// reset - prof only (handlers also gate on PLATFORM env)
	mux.HandleFunc("DELETE /api/reset/students", middleware.RequireProfessor(cfg.HandlerResetStudents, cfg.JWT))       // reset students table
	mux.HandleFunc("DELETE /api/reset/courses", middleware.RequireProfessor(cfg.HandlerResetCourses, cfg.JWT))         // reset courses table
	mux.HandleFunc("DELETE /api/reset/enrollments", middleware.RequireProfessor(cfg.HandlerResetEnrollments, cfg.JWT)) // reset enrollments table
	mux.HandleFunc("DELETE /api/reset/sessions", middleware.RequireProfessor(cfg.HandlerResetSessions, cfg.JWT))       // reset attendance sessions table
	mux.HandleFunc("DELETE /api/reset/records", middleware.RequireProfessor(cfg.HandlerResetRecords, cfg.JWT))         // reset attendance records table
	mux.HandleFunc("DELETE /api/reset/all", middleware.RequireProfessor(cfg.HandlerResetAll, cfg.JWT))                 // reset all tables in correct order
	mux.HandleFunc("DELETE /api/reset/specialties", middleware.RequireProfessor(cfg.HandlerResetSpecialty, cfg.JWT))   // reset table specialties

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
