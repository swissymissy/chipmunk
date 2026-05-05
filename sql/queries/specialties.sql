-- name: CreateSpecialty :one
INSERT INTO specialties (name)
VALUES (?)
RETURNING *;

-- name: ListAllSpecialties :many
SELECT * FROM specialties
ORDER BY name;

-- delete a specialty from the list
-- name: DeleteSpecialty :exec
DELETE FROM specialties WHERE id = ?;

-- reset the specialties table
-- name: ResetSpecialties :exec
DELETE FROM specialties;
