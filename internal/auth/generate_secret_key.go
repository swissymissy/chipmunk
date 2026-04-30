package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// generate random bytes and encode into hex
// for secret key
func GenSecretKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("error generate random bytes: %w")
	}

	// encode to hex string
	hexStr := hex.EncodeToString(bytes)
	return hexStr, nil
}
