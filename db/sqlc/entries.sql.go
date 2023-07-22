// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: entries.sql

package db

import (
	"context"
)

const createEntryRecord = `-- name: CreateEntryRecord :one
INSERT INTO entries (
 account_id,
 amount
) VALUES (
  $1, $2
)
RETURNING id, account_id, amount, created_at
`

type CreateEntryRecordParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntryRecord(ctx context.Context, arg CreateEntryRecordParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntryRecord, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getEntryRecord = `-- name: GetEntryRecord :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEntryRecord(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntryRecord, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEntryRecords = `-- name: ListEntryRecords :many
SELECT id, account_id, amount, created_at FROM entries
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListEntryRecordsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEntryRecords(ctx context.Context, arg ListEntryRecordsParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntryRecords, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
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

const listEntryRecordsByAccountID = `-- name: ListEntryRecordsByAccountID :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListEntryRecordsByAccountIDParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntryRecordsByAccountID(ctx context.Context, arg ListEntryRecordsByAccountIDParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntryRecordsByAccountID, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
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
