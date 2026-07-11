package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")

	// no header
	if header == "" {
		return "", errors.New("invalid header")
	}
	// header start with "bearer"
	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("invalid header")
	}
	// strip "Bearer " to get token
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" {
		return "", errors.New("invalid token")
	}

	return token, nil
}
