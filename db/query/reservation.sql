-- name: GetReservation :one
SELECT * FROM reservations
WHERE id = $1 LIMIT 1;

-- name: GetOptimizedReservation :one
SELECT * FROM reservations
WHERE reservations.start_time = $1 and 
reservations.booked = $2 and 
reservations.table_size = (
  SELECT MIN(table_size)
    FROM (SELECT * FROM reservations 
      WHERE reservations.table_size >= $3 AND reservations.booked = $2 AND 
      reservations.start_time = $1) AS subquery
  ) 
LIMIT 1;


-- name: CreateReservation :one
INSERT INTO reservations (
  table_size, start_time, booked, duration
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateReservation :one
UPDATE reservations
SET table_size = $2,
start_time = $3,
booked = $4,
duration = $5
WHERE id = $1
RETURNING *;

-- name: ListReservations :many
SELECT * FROM reservations
ORDER BY start_time
LIMIT $1
OFFSET $2;

-- name: ListAvailableReservations :many
SELECT * FROM reservations
WHERE booked = false
ORDER BY start_time
LIMIT $1
OFFSET $2;

-- name: DeleteReservation :exec
DELETE FROM reservations
WHERE id = $1;
