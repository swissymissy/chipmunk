package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

// hash password
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return hash, nil
}

// check password and hash
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("error matching password and hash: %w", err)
	}
	return match, nil
}
