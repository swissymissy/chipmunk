package middleware

import "context"

// helper for handlers to retrieve the student ID
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}
