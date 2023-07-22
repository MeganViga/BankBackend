-- name: CreateTranferRecord :one
INSERT INTO transfers (
 from_account_id,
 to_account_id,
 amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransferRecord :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransferRecords :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListTransferRecordsByFromAndToAccountID :many
SELECT * FROM transfers
WHERE from_account_id = $1 and to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;


