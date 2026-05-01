package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// generate qr token using hmac-sha256
func GenerateQRToken(sessionID int64, secretKey string) (string, error) {
	// get current time in integer
	now := time.Now().Unix()

	// payload string = sessionID:timestamp
	payload := fmt.Sprintf("%d:%d", sessionID, now)

	mac := hmac.New(sha256.New, []byte(secretKey)) // create new hmac instance using sha256 and serectKey
	mac.Write([]byte(payload))                     // write data to be signed
	expectedMAC := mac.Sum(nil)
	signature := hex.EncodeToString(expectedMAC) // encode the result hmac to hex string

	// encode payload to base64
	payload64 := base64.URLEncoding.EncodeToString([]byte(payload))

	// token = payload64.signature
	token := payload64 + "." + signature
	return token, nil
}

// decode the token
func ValidateQRToken(token string, secretKey string) (bool, error) {
	// split token into payload and signature
	payload, signature, err := SplitToken(token)
	if err != nil {
		return false, fmt.Errorf("%s", err)
	}

	// create an hmac from secretKey
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)

	// validate
	match := hmac.Equal(signature, expectedMAC)
	if !match {
		return false, fmt.Errorf("invalid token signature")
	}

	payloadStr := string(payload)
	_, timestamp, err := ExtractSessionIDTimestamp(payloadStr)
	if err != nil {
		return false, fmt.Errorf("%s", err)
	}

	// 20s grace period ( 15s + 5s network buffer)
	if time.Now().Unix()-timestamp > 20 {
		return false, fmt.Errorf("token expired")
	}

	return true, nil
}


// helper functions
func SplitToken(token string) ([]byte, []byte, error) {
	// split token into payload and signature to verify
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return nil, nil, fmt.Errorf("malformed token")
	}
	payload64 := split[0]
	signature := split[1]

	// decode signature to get hmac
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding hex string")
	}

	// decode payload64
	payload, err := base64.URLEncoding.DecodeString(payload64)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding base64 payload")
	}

	return payload, sig, nil
}