-- +goose Up
CREATE TABLE courses (
    id TEXT PRIMARY KEY,
    course_name TEXT NOT NULL,
    section_date TEXT NOT NULL,
    start_time TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE courses;