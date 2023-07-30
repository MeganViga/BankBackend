package db

import (
	"context"
	"testing"

	"github.com/MeganViga/BankBackend/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: int64(util.RandomNumber(1,100)),
	}
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n ;i ++{
		go func(){
			result, err :=testStore.TransferTx(context.Background(), arg)
			errs <- err
			results <- result
		}()
		
	}

	for i := 0; i < n;i++{
		err := <- errs
		require.NoError(t, err)
		result := <- results
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
		require.Equal(t, arg.Amount, transfer.Amount)
		_, err = testStore.GetTransferRecord(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// from_account := result.FromAccount
		// require.NotEmpty(t, from_account)
		// require.Equal(t, arg.FromAccountID, from_account.ID)

		// to_account := result.ToAccount
		// require.NotEmpty(t, to_account)
		// require.Equal(t, arg.ToAccountID, to_account.ID)

		from_entry := result.FromEntry
		require.NotEmpty(t, from_entry)
		require.Equal(t, arg.FromAccountID,from_entry.AccountID)
		require.Equal(t, arg.Amount,-from_entry.Amount)
		_, err = testStore.GetEntryRecord(context.Background(), from_entry.ID)
		require.NoError(t, err)
		require.NotZero(t, from_entry.ID)
		require.NotZero(t, from_entry.CreatedAt)

		to_entry := result.ToEntry
		require.NotEmpty(t, to_entry)
		require.Equal(t, arg.ToAccountID,to_entry.AccountID)
		require.Equal(t, arg.Amount,+to_entry.Amount)
		_, err = testStore.GetEntryRecord(context.Background(), to_entry.ID)
		require.NoError(t, err)
		require.NotZero(t, to_entry.ID)
		require.NotZero(t, to_entry.CreatedAt)

		//TODO: check account's balance
		
	}
	
}