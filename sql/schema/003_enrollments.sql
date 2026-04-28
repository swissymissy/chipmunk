-- +goose Up
CREATE TABLE enrollments (
    student_id TEXT NOT NULL,
    course_id TEXT NOT NULL,
    PRIMARY KEY (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE enrollments;