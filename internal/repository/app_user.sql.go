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
	Username  string
	Email     string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
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
SET email = $2,
    first_name = $3,
    last_name = $4,
    is_active = $5
WHERE id = $1
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type UpdateAppUserParams struct {
	ID        uuid.UUID
	Email     string
	FirstName sql.NullString
	LastName  sql.NullString
	IsActive  bool
}

func (q *Queries) UpdateAppUser(ctx context.Context, arg UpdateAppUserParams) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, updateAppUser,
		arg.ID,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.IsActive,
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

const updateAppUserLastLogin = `-- name: UpdateAppUserLastLogin :one
UPDATE app_user
SET last_login = $2
WHERE id = $1
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type UpdateAppUserLastLoginParams struct {
	ID        uuid.UUID
	LastLogin sql.NullTime
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
SET password = $2
WHERE id = $1
    RETURNING id, username, email, password, last_login, first_name, last_name, is_staff, is_active, date_joined
`

type UpdateAppUserPasswordParams struct {
	ID       uuid.UUID
	Password string
}

func (q *Queries) UpdateAppUserPassword(ctx context.Context, arg UpdateAppUserPasswordParams) (AppUser, error) {
	row := q.db.QueryRowContext(ctx, updateAppUserPassword, arg.ID, arg.Password)
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
