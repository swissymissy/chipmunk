-- name: CreateStudent :one
INSERT INTO students (id, student_id, email, first_name, last_name, specialty)
VALUES (?,?,?,?,?,?)
RETURNING *;

-- name: GetByID :one
SELECT * FROM students
WHERE id = ?;

-- name: GetByStudentID :one
SELECT * FROM students
WHERE student_id = ?;

-- name: GetByStudentEmail :one
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

