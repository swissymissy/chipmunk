-- name: NewEnrollment :one
INSERT INTO enrollments (student_id, course_id)
VALUES (?,?)
RETURNING *;

-- name: GetEnrollmentsByStudent :many
SELECT c.id, c.course_name, c.section_date, c.start_time 
FROM courses c 
JOIN enrollments e ON c.id = e.course_id
WHERE e.student_id = ?;

-- show list of all students that have registered for a course
-- name: GetAllStudentsByCourse :many
SELECT s.id, s.student_id, s.email, s.first_name, s.last_name, s.specialty 
FROM students s 
JOIN enrollments e ON s.id = e.student_id
WHERE e.course_id = ?;

-- check if a student enrolls in a specific course
-- name: IsEnrolled :one
SELECT EXISTS(
    SELECT 1 FROM enrollments
    WHERE student_id = ? AND course_id = ?
) AS enrolled;

-- name: ResetEnrollment :exec
DELETE FROM enrollments;

-- let student remove a course from their list
-- name: RemoveACourse :exec
DELETE FROM enrollments WHERE student_id =? AND course_id = ?;

