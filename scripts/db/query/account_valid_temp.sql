-- name: CreateAccountValid :execresult
INSERT INTO account_valid_temp (account_id, code, expired_at)
VALUES (?, ?, ?);

-- name: GetAccountValid :one
SELECT *
FROM account_valid_temp
WHERE code = ?
  AND expired_at > NOW() LIMIT 1;


-- name: DeleteAccountValidByAccountID :exec
DELETE
FROM account_valid_temp
WHERE account_id = ?;