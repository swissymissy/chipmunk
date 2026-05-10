-- find suspicious identical fingerprint for different accounts.
-- returns students whose fingerprint matches at least 1 other student
-- in the same session, regardless of current status — so flag history
-- stays visible after the prof marks a cheater absent.
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
    AND a1.device_fingerprint IS NOT NULL
    AND a1.device_fingerprint != ''
    AND EXISTS (
        SELECT 1 FROM attendance_records a2
        WHERE a2.session_id = a1.session_id
            AND a2.device_fingerprint = a1.device_fingerprint
            AND a2.student_id != a1.student_id
    )
ORDER BY a1.device_fingerprint, a1.check_in_at;

