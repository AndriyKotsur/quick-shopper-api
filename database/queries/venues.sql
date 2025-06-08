-- name: GetVenueByID :one
SELECT * FROM venues WHERE id = $1;

-- name: CreateVenue :one
INSERT INTO venues (id, created_at, updated_at, name, location, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateVenu :one
UPDATE venues
SET updated_at = $2,
name = $3,
location = $4,
type = $5
WHERE id = $1
RETURNING *;

-- name: DeleteVenue :exec
DELETE FROM venues WHERE id = $1;