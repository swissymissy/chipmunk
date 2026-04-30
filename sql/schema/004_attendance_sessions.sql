-- +goose Up
CREATE TABLE attendance_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    course_id TEXT NOT NULL,
    session_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'closed')),
    secret_key TEXT NOT NULL,
    classroom_lat REAL,
    classroom_lng REAL,
    radius_meters INTEGER DEFAULT 50,
    started_at TEXT NOT NULL DEFAULT (datetime('now')),
    ended_at TEXT,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE attendance_sessions;