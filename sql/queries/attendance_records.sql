-- name: CreateRecords :exec
INSERT INTO attendance_records (session_id, student_id)
SELECT ?, e.student_id
FROM enrollments e
WHERE e.course_id = ?;

-- name: StudentCheckIn :one
UPDATE attendance_records
SET status = 'present' , check_in_at = datetime('now'), student_lat = ?, student_lng = ?
WHERE session_id = ? AND student_id = ?
RETURNING *;

-- name: GetRecordBySession :many
SELECT r.* , s.first_name, s.last_name, s.student_id
FROM attendance_records r
JOIN students s ON r.student_id = s.id
WHERE r.session_id = ?;


-- name: UpdateCheckIn :one
UPDATE attendance_records
SET status = 'present', check_in_at = datetime('now')
WHERE student_id = ? AND session_id = ?
RETURNING *;

-- name: ResetAttendanceRecords :exec
DELETE FROM attendance_records;
