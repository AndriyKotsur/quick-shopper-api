-- name: GetVenues :many
SELECT * FROM venues;

-- name: GetVenueByID :one
SELECT * FROM venues WHERE id = $1;

-- name: CreateVenue :one
INSERT INTO venues (id, created_at, updated_at, user_id, name, location, type)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateVenue :one
UPDATE venues
SET updated_at = $2,
user_id = $3,
name = $4,
location = $5,
type = $6
WHERE id = $1
RETURNING *;

-- name: DeleteVenue :exec
DELETE FROM venues WHERE id = $1;