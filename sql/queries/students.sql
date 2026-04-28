-- name: CreateStudent :one
INSERT INTO students (id, student_id, email, password_hash, first_name, last_name, verified, specialty)
VALUES (?,?,?,?,?,?, 1 ,?)
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

