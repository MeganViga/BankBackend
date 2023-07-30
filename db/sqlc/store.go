package db

import (
	"context"
	"database/sql"
	"fmt"
)




type Store struct{
	*Queries
	db *sql.DB
}


func NewStore(db *sql.DB)*Store{
	return &Store{
		db: db,
		Queries: New(db),
	}
}

func (s *Store)execTx(ctx context.Context, fn func(*Queries)error)error{
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil{
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil{
		if rbErr := tx.Rollback(); rbErr != nil{
			return fmt.Errorf("tx err: %v rb err: %v", err, rbErr)
		}
	}
	return tx.Commit()
}
type TransferTxParams struct{
	FromAccountID int64
	ToAccountID int64
	Amount		int64
}
type TransferTxResult struct{
	Transfer Transfer
	FromAccount Account
	ToAccount Account
	FromEntry Entry
	ToEntry Entry
}
func (s *Store)TransferTx(ctx context.Context, arg TransferTxParams)(TransferTxResult, error){
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTranferRecord(ctx, CreateTranferRecordParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		result.FromEntry, err = q.CreateEntryRecord(ctx,CreateEntryRecordParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}
		result.ToEntry, err = q.CreateEntryRecord(ctx,CreateEntryRecordParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		// err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID: arg.FromAccountID,
		// 	Balance: -arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		// err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID: arg.ToAccountID,
		// 	Balance: arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		return nil
	})
	return result, err
}