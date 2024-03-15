-- name: CreateAccount :execresult
INSERT INTO accounts (email, password, name, is_valid)
VALUES (?, ?, ?, ?);

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountByEmail :one
SELECT *
FROM accounts
WHERE email = ? LIMIT 1;

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY id LIMIT ?
OFFSET ?;

-- name: UpdateAccount :execresult
UPDATE accounts
SET password   = ?,
    name       = ?,
    is_valid   = ?,
    updated_at = ?
WHERE id = ?;

-- name: UpdateAccountValid :execresult
UPDATE accounts
SET is_valid   = ?,
    verified_at = ?
WHERE id = ?;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = ?;