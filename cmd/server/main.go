package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/swissymissy/chipmunk/internal/handlers"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	dbURL := os.Getenv("DB_URL")
	// open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to database: %s\n", err)
		return
	}
	// query
	dbQuery := database.New(db)
	log.Print("Database connected")

	platform := os.Getenv("PLATFORM")

	// server config
	cfg := &handlers.ApiConfig{
		Port:     port,
		DB:       dbQuery,
		Platform: platform,
	}
}
