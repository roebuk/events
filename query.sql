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
  slug,
  year)
VALUES ($1, $2, $3, $4)
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


-- Auth Credentials Queries

-- name: CreateAuthCredentials :one
INSERT INTO auth_credentials (
    user_id,
    password_hash
) VALUES ($1, $2)
RETURNING *;

-- name: GetAuthCredentialsByUserID :one
SELECT * FROM auth_credentials
WHERE user_id = $1
AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
AND deleted_at IS NULL
LIMIT 1;

-- name: GetAuthCredentialsByEmail :one
SELECT ac.* FROM auth_credentials ac
INNER JOIN users u ON ac.user_id = u.id
WHERE u.email = $1
AND ac.deleted_at IS NULL
AND u.deleted_at IS NULL
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE auth_credentials
SET last_login_at = NOW(),
    failed_login_attempts = 0,
    locked_until = NULL
WHERE user_id = $1;

-- name: IncrementFailedLoginAttempts :exec
UPDATE auth_credentials
SET failed_login_attempts = failed_login_attempts + 1
WHERE user_id = $1;

-- name: LockAccount :exec
UPDATE auth_credentials
SET locked_until = $2
WHERE user_id = $1;

-- name: IsAccountLocked :one
SELECT
    CASE
        WHEN locked_until IS NULL THEN false
        WHEN locked_until > NOW() THEN true
        ELSE false
    END as is_locked
FROM auth_credentials
WHERE user_id = $1;

-- name: VerifyEmail :exec
UPDATE auth_credentials
SET email_verified_at = NOW()
WHERE user_id = $1;


-- name: CreateRace :one
INSERT INTO races (
  event_id,
  name,
  slug,
  registration_open_date,
  registration_close_date,
  max_capacity,
  price_units,
  currency)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ListRacesByEvent :many
SELECT * FROM races
WHERE event_id = $1
AND deleted_at IS NULL
ORDER BY name;
