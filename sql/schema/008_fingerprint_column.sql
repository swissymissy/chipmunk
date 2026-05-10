-- +goose Up
ALTER TABLE students ADD COLUMN registered_fingerprint TEXT;
ALTER TABLE attendance_records ADD COLUMN device_fingerprint TEXT;

-- +goose Down
ALTER TABLE students DROP COLUMN registered_fingerprint TEXT;
ALTER TABLE attendance_records DROP COLUMN device_fingerprint TEXT;

