-- name: GetEvent :one
SELECT * from events
WHERE id = $1 LIMIT 1;

-- name: ListEvents :many
SELECT * from events
ORDER BY name;

-- name: CreateEvent :one
INSERT INTO events (name, bio)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateEvent :exec
UPDATE events
SET name = $2,
    bio = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;
