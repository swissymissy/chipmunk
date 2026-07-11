package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// create new access token for when student login
func MakeJWT(studentID string, serverSecretToken string) (string, error) {

	// create a new registered claim
	claim := jwt.RegisteredClaims{
		Issuer:    "chipmunk-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(15 * time.Minute)),
		Subject:   studentID,
	}

	// create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// sign the token wuth server secret key
	signedKey := []byte(serverSecretToken)
	signedToken, err := token.SignedString(signedKey)
	if err != nil {
		return "", fmt.Errorf("cannot sign token: %w", err)
	}
	return signedToken, nil
}

// check token
func ValidateJWT(tokenString, serverSecretToken string) (string, error) {
	// create new empty claim struct to be filled
	claim := &jwt.RegisteredClaims{}

	// pass a pointer to that struct so the library can modify it
	_, err := jwt.ParseWithClaims(
		tokenString,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(serverSecretToken), nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("token is expired or bad signature: %w", err)
	}

	return claim.Subject, nil
}
