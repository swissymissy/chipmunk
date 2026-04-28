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

	// professor
	mux.HandleFunc("POST /api/courses", cfg.HandleCreateCourse) // professor create new course

	// students 
	mux.HandleFunc("GET /api/courses", cfg.HandlerGetAllCourses) // list courses to let students pick
	mux.HandleFunc("POST /api/auth/login", middleware.AuthRequired(cfg.HandlerStudentLogin, cfg.JWT))
	mux.HandleFunc("POST /api/auth/register", cfg.HandlerStudentRegister)
	mux.HandleFunc("POST /api/enrollment", middleware.AuthRequired(cfg.HandlerEnrollment, cfg.JWT))

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

	log.Println("Shuttign down server...")

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := chipmunkServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error. Forced shutdown: %s\n", err)
	}
	log.Println("Graceful shutdown complete")
}
