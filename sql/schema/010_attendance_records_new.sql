-- +goose Up
-- create new table for attendance record
-- add 'late' status in for CHECK + note column
CREATE TABLE attendance_records_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    student_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'absent' CHECK(status IN('present', 'absent', 'late')),
    check_in_at TEXT,
    student_lat REAL,
    student_lng REAL,
    accuracy REAL,
    device_fingerprint TEXT,
    note TEXT,
    FOREIGN KEY (session_id) REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    UNIQUE (session_id, student_id)
);

-- copy old exist data from current attendance_records table to attendance_records_new table
INSERT INTO attendance_records_new (
    id, session_id, student_id, status,
    check_in_at, student_lat, student_lng, accuracy,
    device_fingerprint )
SELECT id, session_id, student_id, status, check_in_at, student_lat, student_lng, accuracy, device_fingerprint
FROM attendance_records;

-- drop old attendance_records table
DROP TABLE attendance_records;

-- rename attendance_records_new to attendance_records
ALTER TABLE attendance_records_new RENAME TO attendance_records;

-- +goose Down
-- restore the schema how it was before
CREATE TABLE attendance_records_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    student_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'absent' CHECK(status IN('present', 'absent')),
    check_in_at TEXT,
    student_lat REAL,
    student_lng REAL,
    accuracy REAL,
    device_fingerprint TEXT,
    FOREIGN KEY (session_id) REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    UNIQUE (session_id, student_id)
);

-- copy data from current attendace_records table to new created attendance_records_old table
-- old CHECK does not have 'late', so we count 'late' as 'absent'
INSERT INTO attendance_records_old (id, session_id, student_id, status, check_in_at, student_lat, student_lng, accuracy, device_fingerprint)
SELECT id, session_id, student_id,
    CASE WHEN status = 'late' THEN 'absent' ELSE status END,
    check_in_at, student_lat, student_lng, accuracy, device_fingerprint
FROM attendace_records;
    
DROP TABLE attendance_records;
ALTER TABLE attendance_records_old RENAME TO attendace_records;