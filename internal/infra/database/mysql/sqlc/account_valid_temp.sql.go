// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: account_valid_temp.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createAccountValid = `-- name: CreateAccountValid :execresult
INSERT INTO account_valid_temp (account_id, code, expired_at)
VALUES (?, ?, ?)
`

type CreateAccountValidParams struct {
	AccountID int64     `json:"account_id"`
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (q *Queries) CreateAccountValid(ctx context.Context, arg CreateAccountValidParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createAccountValid, arg.AccountID, arg.Code, arg.ExpiredAt)
}

const deleteAccountValidByAccountID = `-- name: DeleteAccountValidByAccountID :exec
DELETE
FROM account_valid_temp
WHERE account_id = ?
`

func (q *Queries) DeleteAccountValidByAccountID(ctx context.Context, accountID int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccountValidByAccountID, accountID)
	return err
}

const getAccountValid = `-- name: GetAccountValid :one
SELECT id, account_id, code, created_at, expired_at
FROM account_valid_temp
WHERE code = ?
  AND expired_at > NOW() LIMIT 1
`

func (q *Queries) GetAccountValid(ctx context.Context, code string) (AccountValidTemp, error) {
	row := q.db.QueryRowContext(ctx, getAccountValid, code)
	var i AccountValidTemp
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Code,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}
