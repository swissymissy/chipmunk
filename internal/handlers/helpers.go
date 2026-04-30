package handlers

import (
	"database/sql"
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
