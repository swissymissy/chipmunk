-- find suspicious identical fingerprint for different accounts
-- returns present students whose fingerprint matches
--  at least 1 other student in same session
-- name: GetFlaggedFingerprints :many
SELECT 
    a1.device_fingerprint,
    a1.student_id,
    s.student_id AS school_id,
    s.first_name,
    s.last_name,
    a1.check_in_at
FROM attendance_records a1
JOIN students s ON s.id = a1.student_id
WHERE a1.session_id = ?
    AND a1.status = 'present'
    AND a1.device_fingerprint IS NOT NULL
    AND a1.device_fingerprint != ''
    AND EXISTS (
        SELECT 1 FROM attendance_records a2
        WHERE a2.session_id = a1.session_id
            AND a2.device_fingerprint = a1.device_fingerprint
            AND a2.student_id != a1.student_id
            AND a2.status = 'present'
    )
ORDER BY a1.device_fingerprint, a1.check_in_at;

