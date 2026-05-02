package auth

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// helper functions to help process qr token
// split token into payload and signature ( token = payload.signature)
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

// function that extracts session ID and timestamp from the payload
// payload = sessionID:timestamp
func ExtractSessionIDTimestamp(payload string) (int64, int64, error) {
	// split payload into parts
	split := strings.Split(payload, ":")
	if len(split) != 2 {
		return -1, -1, fmt.Errorf("malformed payload")
	}

	// convert to int64
	sessionID, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("error converting string to int")
	}
	timestamp, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("error parsing timestamp")
	}

	return sessionID, timestamp, nil
}
