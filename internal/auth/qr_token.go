package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
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

// validate the qr token
func ValidateQRToken(token string, secretKey string) bool {
	// split token into payload and signature
	payload, signature, err := SplitToken(token)
	if err != nil {
		return false
	}

	// create an hmac from secretKey
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)

	// validate
	match := hmac.Equal(signature, expectedMAC)
	if !match {
		return false
	}

	payloadStr := string(payload)
	_, timestamp, err := ExtractSessionIDTimestamp(payloadStr)
	if err != nil {
		return false
	}

	// 20s grace period ( 15s + 5s network buffer)
	if time.Now().Unix()-timestamp > 20 {
		return false
	}

	return true
}
