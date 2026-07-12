package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// helper function to wrap string to NullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ToNullFloat(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

// convert a SQLite "YYYY-MM-DD HH:MM:SS" UTC timestamp into the server's local time.
// returns the input unchanged if empty or unparseable.
func LocalizeSQLiteTime(s string) string {
	if s == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return s
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

// check student's email
func EmailCheck(s string) (string, error) {
	str := strings.TrimSpace(s)
	if str == "" {
		return "", fmt.Errorf("email can't be empty string")
	}
	if !strings.Contains(str, "@") {
		return "", fmt.Errorf("malformed email")
	}
	return str, nil
}

// check student's name
func NameCheck(s string) (string, error) {
	str := strings.TrimSpace(s)
	if str == "" {
		return "", fmt.Errorf("name can't be empty")
	}
	return str, nil
}

// check student's school ID input
func SchoolIDCheck(s string) (string, error) {
	schoolID := strings.TrimSpace(strings.ToUpper(s))
	if schoolID == "" {
		return "", fmt.Errorf("school ID can't be empty")
	}
	if len(schoolID) != 9 {
		return "", fmt.Errorf("school ID too long")
	}
	for _, r := range schoolID[1:] {
		if r < '0' || r > '9' {
			return "", fmt.Errorf("school ID must contains digit after U. Example: U12345678")
		}
	}
	return schoolID, nil
}
