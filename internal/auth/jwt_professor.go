package auth

import (
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// make jwt token for professor (12h)
func MakeProfessorJWT(serverSecretToken string) (string, error) {

	// create a new registered claim
	claim := jwt.RegisteredClaims{
		Issuer:    "chipmunk-admin",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(12 * time.Hour)),
		Subject:   "professor",
	}

	// create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// sign the token wuth server secret key
	signedKey := []byte(serverSecretToken)
	signedToken, err := token.SignedString(signedKey)
	if err != nil {
		return "", fmt.Errorf("cannot sign professor token: %w", err)
	}
	return signedToken, nil
}