// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package profiledb

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const create2FAInit = `-- name: Create2FAInit :one
INSERT INTO login_2fa (login_uuid, two_factor_secret, two_factor_enabled, two_factor_backup_codes)
VALUES ($1, $2, FALSE, $3::TEXT[])
RETURNING uuid
`

type Create2FAInitParams struct {
	LoginUuid       uuid.UUID   `json:"login_uuid"`
	TwoFactorSecret pgtype.Text `json:"two_factor_secret"`
	Column3         []string    `json:"column_3"`
}

func (q *Queries) Create2FAInit(ctx context.Context, arg Create2FAInitParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, create2FAInit, arg.LoginUuid, arg.TwoFactorSecret, arg.Column3)
	var uuid uuid.UUID
	err := row.Scan(&uuid)
	return uuid, err
}

const disable2FA = `-- name: Disable2FA :exec
UPDATE users
SET two_factor_secret = NULL,
    two_factor_enabled = FALSE,
    two_factor_backup_codes = NULL,
    last_modified_at = NOW()
WHERE uuid = $1
`

func (q *Queries) Disable2FA(ctx context.Context, argUuid uuid.UUID) error {
	_, err := q.db.Exec(ctx, disable2FA, argUuid)
	return err
}

const disable2FAByLoginUuid = `-- name: Disable2FAByLoginUuid :exec
UPDATE login_2fa
SET deleted_at = now() AT TIME ZONE 'utc'
WHERE login_uuid = $1
AND deleted_at IS NULL
`

func (q *Queries) Disable2FAByLoginUuid(ctx context.Context, loginUuid uuid.UUID) error {
	_, err := q.db.Exec(ctx, disable2FAByLoginUuid, loginUuid)
	return err
}

const enable2FA = `-- name: Enable2FA :exec
UPDATE users
SET two_factor_secret = $1::text,
    two_factor_enabled = TRUE,
    two_factor_backup_codes = $2::text[],
    last_modified_at = NOW()
WHERE uuid = $3
`

type Enable2FAParams struct {
	Column1 string    `json:"column_1"`
	Column2 []string  `json:"column_2"`
	Uuid    uuid.UUID `json:"uuid"`
}

func (q *Queries) Enable2FA(ctx context.Context, arg Enable2FAParams) error {
	_, err := q.db.Exec(ctx, enable2FA, arg.Column1, arg.Column2, arg.Uuid)
	return err
}

const enable2FAByLoginUuid = `-- name: Enable2FAByLoginUuid :exec
UPDATE login_2fa
SET two_factor_secret = $1::text,
    two_factor_enabled = TRUE,
    two_factor_backup_codes = $2::text[],
    last_modified_at = NOW()
WHERE login_uuid = $3
AND deleted_at IS NULL
`

type Enable2FAByLoginUuidParams struct {
	Column1   string    `json:"column_1"`
	Column2   []string  `json:"column_2"`
	LoginUuid uuid.UUID `json:"login_uuid"`
}

func (q *Queries) Enable2FAByLoginUuid(ctx context.Context, arg Enable2FAByLoginUuidParams) error {
	_, err := q.db.Exec(ctx, enable2FAByLoginUuid, arg.Column1, arg.Column2, arg.LoginUuid)
	return err
}

const findUserByUsername = `-- name: FindUserByUsername :many
SELECT uuid, username, email, password, two_factor_secret, two_factor_enabled, two_factor_backup_codes, created_at, last_modified_at
FROM users
WHERE username = $1
`

type FindUserByUsernameRow struct {
	Uuid                 uuid.UUID      `json:"uuid"`
	Username             sql.NullString `json:"username"`
	Email                string         `json:"email"`
	Password             []byte         `json:"password"`
	TwoFactorSecret      pgtype.Text    `json:"two_factor_secret"`
	TwoFactorEnabled     pgtype.Bool    `json:"two_factor_enabled"`
	TwoFactorBackupCodes []string       `json:"two_factor_backup_codes"`
	CreatedAt            time.Time      `json:"created_at"`
	LastModifiedAt       time.Time      `json:"last_modified_at"`
}

func (q *Queries) FindUserByUsername(ctx context.Context, username sql.NullString) ([]FindUserByUsernameRow, error) {
	rows, err := q.db.Query(ctx, findUserByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindUserByUsernameRow
	for rows.Next() {
		var i FindUserByUsernameRow
		if err := rows.Scan(
			&i.Uuid,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.TwoFactorSecret,
			&i.TwoFactorEnabled,
			&i.TwoFactorBackupCodes,
			&i.CreatedAt,
			&i.LastModifiedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUUID = `-- name: GetUserByUUID :one

SELECT uuid, username, email, password, created_at, last_modified_at
FROM users
WHERE uuid = $1
`

type GetUserByUUIDRow struct {
	Uuid           uuid.UUID      `json:"uuid"`
	Username       sql.NullString `json:"username"`
	Email          string         `json:"email"`
	Password       []byte         `json:"password"`
	CreatedAt      time.Time      `json:"created_at"`
	LastModifiedAt time.Time      `json:"last_modified_at"`
}

// Verify current password
func (q *Queries) GetUserByUUID(ctx context.Context, argUuid uuid.UUID) (GetUserByUUIDRow, error) {
	row := q.db.QueryRow(ctx, getUserByUUID, argUuid)
	var i GetUserByUUIDRow
	err := row.Scan(
		&i.Uuid,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.LastModifiedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2,
    last_modified_at = NOW()
WHERE uuid = $1
`

type UpdateUserPasswordParams struct {
	Uuid     uuid.UUID `json:"uuid"`
	Password []byte    `json:"password"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.Exec(ctx, updateUserPassword, arg.Uuid, arg.Password)
	return err
}

const updateUsername = `-- name: UpdateUsername :exec
UPDATE users
SET username = $2,
    last_modified_at = NOW()
WHERE uuid = $1
`

type UpdateUsernameParams struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Username sql.NullString `json:"username"`
}

func (q *Queries) UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error {
	_, err := q.db.Exec(ctx, updateUsername, arg.Uuid, arg.Username)
	return err
}
