package handlers

import "github.com/swissymissy/chipmunk/internal/database"

type ApiConfig struct {
	Port                  string
	DB                    *database.Queries
	Platform              string
	JWT                   string
	BaseURL               string
	ProfessorPasswordHash string
}