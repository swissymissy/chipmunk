-- +goose Up
CREATE TABLE attendance_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    student_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'absent' CHECK(status IN('present', 'absent')),
    check_in_at TEXT,
    student_lat REAL,
    student_lng REAL,
    FOREIGN KEY (session_id) REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    UNIQUE (session_id, student_id)
);

-- +goose Down
DROP TABLE attendance_records;
