-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: GetTransfers :many
SELECT * FROM transfers
WHERE 
    from_account_id = $1 OR
    to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: GetTransfersByAccount :many
SELECT * FROM transfers
WHERE 
    from_account_id = sqlc.arg(id) OR
    to_account_id = sqlc.arg(id)
ORDER BY id
LIMIT sqlc.arg(size)
OFFSET sqlc.arg(off);