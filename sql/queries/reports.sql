-- semester report: count attendance per student for a course
-- name: GetAttendanceSummaryByCourse :many
SELECT 
    s.student_id,
    s.first_name,
    s.last_name,
    s.specialty,
    COUNT(CASE WHEN r.status='present' THEN 1 END) AS total_present,
    COUNT(r.id) AS total_sessions,
    ROUND(COUNT(CASE WHEN r.status='present' THEN 1 END)*100.0/COUNT(r.id), 1) AS average
FROM students s 
JOIN attendance_records r ON s.id = r.student_id
JOIN attendance_sessions sess ON r.session_id = sess.id 
WHERE sess.course_id = ?
GROUP BY s.id 
ORDER BY s.last_name, s.first_name;

-- get attendance records by date range
-- name: GetAttendanceSummaryByCourseInDateRange :many
SELECT 
    s.student_id,
    s.first_name,
    s.last_name,
    s.specialty,
    COUNT(CASE WHEN r.status='present' THEN 1 END) AS total_present,
    COUNT(r.id) AS total_sessions,
    ROUND(COUNT(CASE WHEN r.status='present' THEN 1 END)*100.0/COUNT(r.id), 1) AS average
FROM students s 
JOIN attendance_records r ON r.student_id = s.id 
JOIN attendance_sessions sess ON r.session_id = sess.id 
WHERE sess.course_id = ? AND sess.session_date >= ? AND sess.session_date <= ?
GROUP BY s.id 
ORDER BY s.last_name, s.first_name;

-- daily report: all records for sessions on a specific date
-- name: GetAttendanceByDate :many
SELECT
    sess.session_date,
    c.course_name,
    c.start_time,
    s.student_id,
    s.first_name,
    s.last_name,
    r.status,
    r.check_in_at
FROM attendance_records r 
JOIN students s ON r.student_id = s.id 
JOIN attendance_sessions sess ON r.session_id = sess.id 
JOIN courses c ON sess.course_id = c.id 
WHERE sess.session_date = ?
ORDER BY c.start_time, s.last_name, s.first_name;

