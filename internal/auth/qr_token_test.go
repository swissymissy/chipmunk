package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"
)

const testSecret = "test-secret-key"

// build a token with an arbitrary timestamp
func buildToken(sessionID int64, timestamp int64, secretKey string) string {
	payload := fmt.Sprintf("%d:%d", sessionID, timestamp)
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))
	payload64 := base64.URLEncoding.EncodeToString([]byte(payload))
	return payload64 + "." + signature
}

func TestGenerateQRToken(t *testing.T) {
	var sessionID int64 = 42
	token, err := GenerateQRToken(sessionID, testSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		t.Fatalf("expected token to have 2 parts separated by '.', got %d", len(parts))
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("expected payload to be valid base64, got error: %v", err)
	}

	payload := string(payloadBytes)
	if !strings.HasPrefix(payload, fmt.Sprintf("%d:", sessionID)) {
		t.Fatalf("expected payload to start with %d:, got %q", sessionID, payload)
	}
}

func TestSplitToken_Valid(t *testing.T) {
	token, err := GenerateQRToken(42, testSecret)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	payload, signature, err := SplitToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(payload) == 0 {
		t.Fatal("expected non-empty payload")
	}
	if len(signature) == 0 {
		t.Fatal("expected non-empty signature")
	}
}

func TestSplitToken_NoDot(t *testing.T) {
	if _, _, err := SplitToken("nodothereatall"); err == nil {
		t.Fatal("expected error for token without '.'")
	}
}

func TestSplitToken_Empty(t *testing.T) {
	if _, _, err := SplitToken(""); err == nil {
		t.Fatal("expected error for empty string")
	}
}

func TestSplitToken_InvalidBase64(t *testing.T) {
	// valid hex signature, but '!' is not a valid base64 char
	if _, _, err := SplitToken("!!!!.deadbeef"); err == nil {
		t.Fatal("expected error for invalid base64 payload")
	}
}

func TestSplitToken_InvalidHex(t *testing.T) {
	// valid base64 payload, but signature has non-hex chars
	validPayload := base64.URLEncoding.EncodeToString([]byte("42:1714500000"))
	if _, _, err := SplitToken(validPayload + ".not-hex-zz"); err == nil {
		t.Fatal("expected error for invalid hex signature")
	}
}

func TestExtractSessionIDTimestamp_Valid(t *testing.T) {
	sid, ts, err := ExtractSessionIDTimestamp("42:1714500000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sid != 42 {
		t.Errorf("expected sessionID 42, got %d", sid)
	}
	if ts != 1714500000 {
		t.Errorf("expected timestamp 1714500000, got %d", ts)
	}
}

func TestExtractSessionIDTimestamp_NoColon(t *testing.T) {
	if _, _, err := ExtractSessionIDTimestamp("421714500000"); err == nil {
		t.Fatal("expected error for payload without colon")
	}
}

func TestExtractSessionIDTimestamp_NonNumericSessionID(t *testing.T) {
	if _, _, err := ExtractSessionIDTimestamp("abc:1714500000"); err == nil {
		t.Fatal("expected error for non-numeric sessionID")
	}
}

func TestExtractSessionIDTimestamp_NonNumericTimestamp(t *testing.T) {
	if _, _, err := ExtractSessionIDTimestamp("42:notatime"); err == nil {
		t.Fatal("expected error for non-numeric timestamp")
	}
}

func TestExtractSessionIDTimestamp_Empty(t *testing.T) {
	if _, _, err := ExtractSessionIDTimestamp(""); err == nil {
		t.Fatal("expected error for empty payload")
	}
}

func TestValidateQRToken_SameSecret(t *testing.T) {
	token, err := GenerateQRToken(42, testSecret)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if !ValidateQRToken(token, testSecret) {
		t.Fatal("expected token to validate with same secret key")
	}
}

func TestValidateQRToken_DifferentSecret(t *testing.T) {
	token, err := GenerateQRToken(42, testSecret)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if ValidateQRToken(token, "different-secret") {
		t.Fatal("expected token to fail with different secret key")
	}
}

func TestValidateQRToken_Expired(t *testing.T) {
	// 60s in the past — well outside the 20s grace window
	oldTime := time.Now().Unix() - 60
	token := buildToken(42, oldTime, testSecret)
	if ValidateQRToken(token, testSecret) {
		t.Fatal("expected expired token to fail validation")
	}
}

func TestValidateQRToken_TamperedPayload(t *testing.T) {
	token, err := GenerateQRToken(42, testSecret)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	parts := strings.Split(token, ".")
	// swap payload while keeping the original signature
	tamperedPayload := base64.URLEncoding.EncodeToString([]byte("99:1714500000"))
	tampered := tamperedPayload + "." + parts[1]

	if ValidateQRToken(tampered, testSecret) {
		t.Fatal("expected validation to fail for tampered payload")
	}
}

func TestValidateQRToken_TamperedSignature(t *testing.T) {
	token, err := GenerateQRToken(42, testSecret)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	parts := strings.Split(token, ".")
	// flip the first hex character of the signature
	sig := parts[1]
	var flipped byte = '1'
	if sig[0] == '1' {
		flipped = '0'
	}
	tamperedSig := string(flipped) + sig[1:]
	tampered := parts[0] + "." + tamperedSig

	if ValidateQRToken(tampered, testSecret) {
		t.Fatal("expected validation to fail for tampered signature")
	}
}

func TestValidateQRToken_Malformed(t *testing.T) {
	if ValidateQRToken("not-a-real-token", testSecret) {
		t.Fatal("expected validation to fail for malformed token")
	}
}
