-- +goose Up
ALTER TABLE students 
ADD COLUMN must_change_password INTEGER NOT NULL
DEFAULT 0;

-- +goose Down
ALTER TABLE students 
DROP COLUMN must_change_password;

