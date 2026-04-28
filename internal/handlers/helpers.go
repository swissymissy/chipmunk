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
