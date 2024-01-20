-- name: GetAppUserById :one
SELECT * FROM app_user
WHERE id = $1 LIMIT 1;

-- name: GetAppUserByEmailAddr :one
SELECT * FROM app_user
WHERE email = $1 LIMIT 1;

-- name: ListAppUser :many
SELECT * FROM app_user;

-- name: CreateAppUser :one
INSERT INTO app_user (
    username,
    email,
    first_name,
    last_name,
    is_active
) VALUES (
             $1, $2, $3, $4, $5
         )
    RETURNING *;

-- name: UpdateAppUser :one
UPDATE app_user
SET email = $2,
    first_name = $3,
    last_name = $4,
    is_active = $5,
    last_login = $6
WHERE id = $1
    RETURNING *;
