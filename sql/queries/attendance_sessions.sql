-- name: CreateSession :one
INSERT INTO attendance_sessions (
    course_id,
    session_date,
    secret_key,
    classroom_lat,
    classroom_lng,
    radius_meters
)
VALUES (?,?,?,?,?,?)
RETURNING *;

-- name: CloseSession :one
UPDATE attendance_sessions
SET status = 'closed', ended_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: ReOpenSession :one
UPDATE attendance_sessions
SET status = 'active', ended_at = NULL
WHERE id = ? 
RETURNING *;

-- get active session for a course (to check if one already exists)
-- name: GetActiveSession :one
SELECT * FROM attendance_sessions
WHERE course_id = ? AND status = 'active';

-- get session by ID
-- name: GetSessionByID :one
SELECT * FROM attendance_sessions
WHERE id = ?;

-- list all active sessions in case professor forget to close session
-- name: ListActiveSessions :many
SELECT s.id, s.course_id, s.session_date, s.status, s.started_at, c.course_name
FROM attendance_sessions s 
JOIN courses c ON s.course_id = c.id
WHERE s.status = 'active';

-- name: DeleteSession :exec
DELETE FROM attendance_sessions
WHERE id = ?;

-- name: ResetAttendanceSession :exec
DELETE FROM attendance_sessions; 