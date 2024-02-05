-- name: GetAppUserById :one
SELECT * FROM app_user
WHERE id = $1 LIMIT 1;

-- name: GetAppUserByEmailAddr :one
SELECT * FROM app_user
WHERE email = $1 LIMIT 1;

-- name: GetAppUserByUsername :one
SELECT * FROM app_user
WHERE username = $1 LIMIT 1;

-- name: ListAppUser :many
SELECT * FROM app_user;

-- name: CreateAppUser :one
INSERT INTO app_user (
    username,
    email,
    password,
    first_name,
    last_name
) VALUES (
             $1, $2, $3, $4, $5
         )
    RETURNING *;

-- name: UpdateAppUser :one
UPDATE app_user
SET first_name = coalesce(sqlc.narg('first_name'), first_name),
    last_name = coalesce(sqlc.narg('last_name'), last_name)
WHERE id = sqlc.arg('id')
    RETURNING *;

-- name: UpdateAppUserPassword :one
UPDATE app_user
SET password = $2
WHERE id = $1
    RETURNING *;

-- name: UpdateAppUserLastLogin :one
UPDATE app_user
SET last_login = $2
WHERE id = $1
    RETURNING *;