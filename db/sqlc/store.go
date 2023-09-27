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
		return err
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
var txKey = struct{}{}
func (s *Store)TransferTx(ctx context.Context, arg TransferTxParams)(TransferTxResult, error){
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName,"Create Transfer")
		result.Transfer, err = q.CreateTranferRecord(ctx, CreateTranferRecordParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		fmt.Println(txName,"Create entry 1")
		result.FromEntry, err = q.CreateEntryRecord(ctx,CreateEntryRecordParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}
		fmt.Println(txName,"Create entry 2")
		result.ToEntry, err = q.CreateEntryRecord(ctx,CreateEntryRecordParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		// account1, err := q.GetAccount(ctx, arg.FromAccountID)
		// if err != nil{
		// 	return err
		// }
		// result.FromAccount, err =q.UpdateAccount(ctx, UpdateAccountParams{
		// 	account1.ID,
		// 	account1.Balance - arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		// account2, err := q.GetAccount(ctx, arg.ToAccountID)
		// if err != nil{
		// 	return err
		// }
		// result.ToAccount, err =q.UpdateAccount(ctx, UpdateAccountParams{
		// 	account2.ID,
		// 	account2.Balance + arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		// fmt.Println(txName,"Get Account 1 for update")
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil{
		// 	return err
		// }
		// result.FromAccount, err =q.UpdateAccount(ctx, UpdateAccountParams{
		// 	account1.ID,
		// 	account1.Balance - arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		if arg.FromAccountID < arg.ToAccountID{
			result.FromAccount, result.ToAccount, err =addMoney(ctx, q, arg.FromAccountID, -arg.Amount,arg.ToAccountID,arg.Amount)
			// fmt.Println(txName,"Update Account 1")
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID: arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil{
			// 	return err
			// }


			// fmt.Println(txName,"Update Account 2")
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID: arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil{
			// 	return err
			// }

		}else{
			result.ToAccount, result.FromAccount, err =addMoney(ctx, q, arg.ToAccountID, arg.Amount,arg.FromAccountID,-arg.Amount)
			// fmt.Println(txName,"Update Account 1")
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID: arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil{
			// 	return err
			// }


			// fmt.Println(txName,"Update Account 2")
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID: arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil{
			// 	return err
			// }

		}
		
		// fmt.Println(txName,"Get Account 2 for update")
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil{
		// 	return err
		// }
		// result.ToAccount, err =q.UpdateAccount(ctx, UpdateAccountParams{
		// 	account2.ID,
		// 	account2.Balance + arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		
		return nil
	})
	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64)(account1 Account, account2 Account, err error){
			account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID: accountID1,
				Amount: amount1,
			})
			if err != nil{
				return 
			}
			account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID: accountID2,
				Amount: amount2,
			})
			return
}
