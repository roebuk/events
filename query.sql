-- name: GetEvent :one
SELECT * from events
WHERE slug = $1 LIMIT 1;

-- name: ListEvents :many
SELECT * from events
ORDER BY name;


-- name: CreateEvent :one
INSERT INTO events (
  organisation_id,
  name,
  slug)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateEvent :exec
UPDATE events
SET name = $2,
    slug = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
UPDATE events
SET deleted_at = NOW()
WHERE id = $1;


-- name: GetOrganisation :one
SELECT * from organisations
WHERE id = $1 LIMIT 1;

-- name: CreateOrganisation :one
INSERT INTO organisations (
  name)
VALUES ($1)
RETURNING *;

-- name: UpdateOrganisation :exec
UPDATE organisations
SET name = $2
WHERE id = $1
RETURNING *;


-- name: DeleteOrganisation :exec
UPDATE organisations
SET deleted_at = NOW()
WHERE id = $1;


-- name: GetUser :one
SELECT * from users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  email,
  first_name,
  last_name,
  phone,
  address_line1,
  address_line2,
  city,
  state,
  postal_code,
  country,
  role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET email = $2,
    first_name = $3,
    last_name = $4,
    phone = $5,
    address_line1 = $6,
    address_line2 = $7,
    city = $8,
    state = $9,
    postal_code = $10,
    country = $11,
    role = $12
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;
