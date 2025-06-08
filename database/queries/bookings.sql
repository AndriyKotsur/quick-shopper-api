-- name: GetBookings :many
SELECT * FROM bookings;

-- name: GetBookingByID :one
SELECT * FROM bookings WHERE id = $1;

-- name: CreateBooking :one
INSERT INTO bookings (id, created_at, updated_at, user_id, venue_id, start_time, end_time, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateBooking :one
UPDATE bookings
SET updated_at = $2,
user_id = $3,
venue_id = $4,
start_time = $5,
end_time = $6,
status = $7
WHERE id = $1
RETURNING *;

-- name: DeleteBooking :exec
DELETE FROM bookings WHERE id = $1;