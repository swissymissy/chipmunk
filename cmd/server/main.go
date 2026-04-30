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
	mux.HandleFunc("POST /api/courses", middleware.LocalOnly(cfg.HandleCreateCourse))                   // professor create new course
	mux.HandleFunc("POST /api/sessions/start", middleware.LocalOnly(cfg.HandlerStartSession))           // start a new session
	mux.HandleFunc("POST /api/sessions/close", middleware.LocalOnly(cfg.HandlerCloseSession))           // close an active session
	mux.HandleFunc("GET /api/sessions/reopen", middleware.LocalOnly(cfg.HandlerReopenSession))          // reopen a closed session
	mux.HandleFunc("GET /api/sessions/{id}", middleware.LocalOnly(cfg.HandlerSessionDetail))            // get session details
	mux.HandleFunc("GET /api/roster/{course_id}", middleware.LocalOnly(cfg.HandlerRosters))             // view students enrolled in a course
	mux.HandleFunc("PUT /api/attendance/override", middleware.LocalOnly(cfg.HandlerMarkStudentPresent)) // manually mark a student present

	// students
	mux.HandleFunc("GET /api/courses", cfg.HandlerGetAllCourses) // list courses to let students pick
	mux.HandleFunc("POST /api/auth/login", cfg.HandlerStudentLogin)
	mux.HandleFunc("POST /api/auth/register", cfg.HandlerStudentRegister)

	// students - auth required
	mux.HandleFunc("POST /api/enrollment", middleware.AuthRequired(cfg.HandlerEnrollment, cfg.JWT))
	mux.HandleFunc("POST /api/attendance/checkin", middleware.AuthRequired(cfg.HandlerStudentCheckIn, cfg.JWT)) // students check in

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
