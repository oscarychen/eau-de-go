// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: app_user.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createAppUser = `-- name: CreateAppUser :one
INSERT INTO app_user (
    username,
    email,
    password,
    first_name,
    last_name
) VALUES (
             $1, $2, $3, $4, $5
         )
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type CreateAppUserParams struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (q *Queries) CreateAppUser(ctx context.Context, arg CreateAppUserParams) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, createAppUser,
		arg.Username,
		arg.Email,
		arg.Password,
		arg.FirstName,
		arg.LastName,
	)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}

const getAppUserByEmailAddr = `-- name: GetAppUserByEmailAddr :one
SELECT id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined FROM app_user
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetAppUserByEmailAddr(ctx context.Context, email string) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, getAppUserByEmailAddr, email)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}

const getAppUserById = `-- name: GetAppUserById :one
SELECT id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined FROM app_user
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAppUserById(ctx context.Context, id uuid.UUID) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, getAppUserById, id)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}

const getAppUserByUsername = `-- name: GetAppUserByUsername :one
SELECT id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined FROM app_user
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetAppUserByUsername(ctx context.Context, username string) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, getAppUserByUsername, username)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}

const listAppUser = `-- name: ListAppUser :many
SELECT id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined FROM app_user
`

func (q *Queries) ListAppUser(ctx context.Context) ([]AppUser, error) {
	rows, err := q.db.QueryContext(ctx, listAppUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppUser
	for rows.Next() {
		var i AppUser
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.LastLogin,
			&i.FirstName,
			&i.LastName,
			&i.IsStaff,
			&i.IsActive,
			&i.DateJoined,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAppUser = `-- name: UpdateAppUser :one
UPDATE app_user
SET first_name = coalesce($1, first_name),
    last_name = coalesce($2, last_name)
WHERE id = $3
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type UpdateAppUserParams struct {
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	ID        uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateAppUser(ctx context.Context, arg UpdateAppUserParams) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, updateAppUser, arg.FirstName, arg.LastName, arg.ID)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}

const updateAppUserLastLogin = `-- name: UpdateAppUserLastLogin :one
UPDATE app_user
SET last_login = $2
WHERE id = $1
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type UpdateAppUserLastLoginParams struct {
	ID        uuid.UUID    `json:"id"`
	LastLogin sql.NullTime `json:"last_login"`
}

func (q *Queries) UpdateAppUserLastLogin(ctx context.Context, arg UpdateAppUserLastLoginParams) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, updateAppUserLastLogin, arg.ID, arg.LastLogin)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}

const updateAppUserPassword = `-- name: UpdateAppUserPassword :one
UPDATE app_user
SET password = $1
WHERE id = $2
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type UpdateAppUserPasswordParams struct {
	Password string    `json:"password"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UpdateAppUserPassword(ctx context.Context, arg UpdateAppUserPasswordParams) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, updateAppUserPassword, arg.Password, arg.ID)
	var i AppUser
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.LastLogin,
		&i.FirstName,
		&i.LastName,
		&i.IsStaff,
		&i.IsActive,
		&i.DateJoined,
	)
	return i, err
}
