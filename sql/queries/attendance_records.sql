-- name: CreateRecords :exec
INSERT INTO attendance_records (session_id, student_id)
SELECT ?, e.student_id
FROM enrollments e
WHERE e.course_id = ?;

-- using UPSERT to insert and update the new registered students
-- name: StudentCheckIn :one
INSERT INTO attendance_records (session_id, student_id, status, check_in_at, student_lat, student_lng, accuracy, device_fingerprint)
VALUES (?,?, 'present', datetime('now'), ?, ?, ? ,? )
ON CONFLICT(session_id, student_id) DO UPDATE SET
    status = 'present',
    check_in_at = datetime('now'),
    student_lat = excluded.student_lat,
    student_lng = excluded.student_lng,
    accuracy = excluded.accuracy,
    device_fingerprint = excluded.device_fingerprint
RETURNING *;

-- get list of students in a session to check their status
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

-- flip a student's record back to 'absent' (prof override after a flag review).
-- preserves device_fingerprint and check_in_at so the flag history stays
-- visible in GetFlaggedFingerprints.
-- name: RevertCheckin :one
UPDATE attendance_records
SET status = 'absent'
WHERE student_id = ? AND session_id = ?
RETURNING *;

-- name: ResetAttendanceRecords :exec
DELETE FROM attendance_records;
