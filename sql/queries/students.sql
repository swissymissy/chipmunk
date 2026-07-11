-- name: CreateStudent :one
INSERT INTO students (id, student_id, email, password_hash, first_name, last_name, verified, specialty, registered_fingerprint)
VALUES (?,?,?,?,?,?, 1 ,?, ?)
RETURNING *;

-- name: GetByID :one
SELECT * FROM students
WHERE id = ?;

-- name: GetStudentByID :one
SELECT * FROM students
WHERE student_id = ?;

-- name: GetStudentByEmail :one
SELECT * FROM students
WHERE email = ?;

-- name: UpdateStudentEmail :one
UPDATE students
SET email = ?, updated_at = datetime('now')
WHERE student_id = ?
RETURNING *;

-- name: UpdatePassword :one
UPDATE students
SET password_hash = ?, updated_at = datetime('now'), verified = 1
WHERE student_id = ?
RETURNING *;

-- name: ListAllStudents :many
SELECT * FROM students;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE id = ?;

-- name: ResetStudents :exec
DELETE FROM students;

-- student editing profile feature
-- Search for one student's profile by the UUID

-- name: GetProfileByID :one
SELECT id, student_id, first_name, last_name, specialty
FROM students WHERE id = ?;

-- name: UpdateStudentSchoolID :one
UPDATE students 
SET student_id = ? , updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: UpdateStudentEmailByID :one
UPDATE students 
SET email = ?, updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: UpdateStudentName :one
UPDATE students 
SET first_name = ?, last_name = ? , updated_at = datetime('now')
WHERE id = ?
RETURNING *;


