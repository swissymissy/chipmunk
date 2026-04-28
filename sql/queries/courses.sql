-- name: CreateCourse :one
INSERT INTO courses (id, course_name, section_date, start_time)
VALUES (?,?,?,?)
RETURNING *;

-- name: ListAllCourses :many
SELECT * FROM courses;

-- name: GetCourseByID :one
SELECT * FROM courses WHERE id = ?;

-- name: DeleteCourse :exec
DELETE FROM courses 
WHERE id = ?;

-- name: UpdateCourseName :exec
UPDATE courses 
SET course_name = ? WHERE id = ?;

-- name: UpdateCourseSection :exec
UPDATE courses 
SET section_date = ? WHERE id = ?;

-- name: UpdateCourseStartTime :exec
UPDATE courses
SET start_time = ? WHERE id = ?;

-- name: DeleteAllCourse :exec
DELETE FROM courses; 
