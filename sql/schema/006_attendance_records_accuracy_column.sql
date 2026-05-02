-- +goose Up
ALTER TABLE attendance_records ADD COLUMN accuracy REAL;

-- +goose Down
ALTER TABLE attendance_records DROP COLUMN accuracy;