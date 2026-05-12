package handlers

import (
	"sync"

	"github.com/swissymissy/chipmunk/internal/database"
)

type ApiConfig struct {
	Port                  string
	DB                    *database.Queries
	Platform              string
	JWT                   string
	ProfessorPasswordHash string

	baseURL string
	mu      sync.RWMutex
}

// get new base url
func (cfg *ApiConfig) GetBaseURL() string {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	return cfg.baseURL
}

// update base URL to new generated url from cloudflared tunnel
func (cfg *ApiConfig) SetBaseURL(url string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	cfg.baseURL = url
}
