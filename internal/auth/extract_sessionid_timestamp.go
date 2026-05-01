package auth

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// function that extract session ID and timestamp from the qr token
func ExtractSessionIDTimestamp(token string) (int64, int64, error) {
	split := strings.Split(token, ".")
	if len(split) != 2 {
		return -1, -1, fmt.Errorf("malformed token")
	} 

	// decode payload to bytes
 	payload, err := base64.URLEncoding.DecodeString(split[0])
	if err != nil {
		return -1, -1, fmt.Errorf("malformed payload")
	}
	payloadStr := string(payload) // convert to string

	// split payloay to get sessionID and timestamp (payload = sessionID:timestamp)
	parts := strings.Split(payloadStr, ":")
	if len(parts) != 2 {
		return -1,-1, fmt.Errorf("malformed payload")
	}

	// convert to int64
	sessionID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("error converting string to int")
	}
	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1, -1, fmt.Errorf("error parsing timestamp")
	}
	
	return sessionID, timestamp, nil
}