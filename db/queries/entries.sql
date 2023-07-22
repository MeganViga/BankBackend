-- name: CreateEntryRecord :one
INSERT INTO entries (
 account_id,
 amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntryRecord :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntryRecords :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListEntryRecordsByAccountID :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;


