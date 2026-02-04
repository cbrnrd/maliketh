-- name: GetOperatorByUsername :one
SELECT * FROM operators WHERE username = $1;

-- name: GetOperatorByAuthToken :one
SELECT * FROM operators WHERE auth_token = $1;

-- name: CreateOperator :one
INSERT INTO operators (
    username, public_key, verify_key, login_secret, rmq_queue, role
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateOperatorAuthToken :exec
UPDATE operators
SET auth_token = $2, auth_token_expiry = $3
WHERE username = $1;

-- name: UpdateOperatorLastLogin :exec
UPDATE operators
SET last_login = NOW()
WHERE username = $1;

-- name: RevokeOperator :exec
UPDATE operators
SET revoked = true
WHERE username = $1;

-- name: ListOperators :many
SELECT * FROM operators ORDER BY created_at DESC;

-- name: CountOperators :one
SELECT COUNT(*) FROM operators;
