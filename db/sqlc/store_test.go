package db

import (
	"context"
	"fmt"
	"testing"

	//"github.com/MeganViga/BankBackend/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>>> before update account 1 balance", account1.Balance)
	fmt.Println(">>>> before update account 2 balance", account2.Balance)
	
	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: int64(10),
	}
	n := 5
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n ;i ++{
		txName := fmt.Sprintf("tx %d",i+1)
		go func(){
			//fmt.Println("txKey", txKey)
			ctx := context.WithValue(context.Background(),txKey,txName)
			result, err :=testStore.TransferTx(ctx, arg, )
			errs <- err
			results <- result
		}()
		
	}
	existed := make(map[int64]bool)
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

		//check accounts
		from_account := result.FromAccount
		require.NotEmpty(t, from_account)
		require.Equal(t, arg.FromAccountID, from_account.ID)

		to_account := result.ToAccount
		require.NotEmpty(t, to_account)
		require.Equal(t, arg.ToAccountID, to_account.ID)

		//TODO: check account's balance
		fmt.Println(">> tx:", from_account.Balance, to_account.Balance)
		diff1 := account1.Balance - from_account.Balance
		diff2 := to_account.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % arg.Amount == 0)
		k := diff1 / arg.Amount
		require.True(t, k >= 1 && k <= int64(n))
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the updated final balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>>> after update account 1 balance", updatedAccount1.Balance)
	fmt.Println(">>>> after update account 2 balance", updatedAccount2.Balance)
	diff1 := account1.Balance - updatedAccount1.Balance
	diff2 := updatedAccount2.Balance - account2.Balance
	require.Equal(t, diff1, diff2)
	require.True(t, diff1 > 0)
	require.True(t, diff1 % arg.Amount == 0)
	k := diff1 / arg.Amount
	require.Equal(t, k, int64(n))
	require.Equal(t, account1.Balance - arg.Amount * int64(n),updatedAccount1.Balance)
	require.Equal(t, account2.Balance + arg.Amount * int64(n),updatedAccount2.Balance)
	
}

func TestTransferTxDeadlock(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>>> before update account 1 balance", account1.Balance)
	fmt.Println(">>>> before update account 2 balance", account2.Balance)
	
	n := 10
	errs := make(chan error)
	// results := make(chan TransferTxResult)
	for i := 0; i < n ;i ++{
		txName := fmt.Sprintf("tx %d",i+1)
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i % 2 == 1{
		fromAccountID = account2.ID
		toAccountID = account1.ID
		}
		go func(){
			//fmt.Println("txKey", txKey)
			arg := TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: int64(10),
			}
			ctx := context.WithValue(context.Background(),txKey,txName)
			_, err :=testStore.TransferTx(ctx, arg, )
			errs <- err
			// results <- result
		}()
		
	}
	// existed := make(map[int64]bool)
	for i := 0; i < n;i++{
		err := <- errs
		require.NoError(t, err)
		// result := <- results
		// require.NotEmpty(t, result)
		// transfer := result.Transfer
		// require.NotEmpty(t, transfer)
		// require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		// require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
		// require.Equal(t, arg.Amount, transfer.Amount)
		// _, err = testStore.GetTransferRecord(context.Background(), transfer.ID)
		// require.NoError(t, err)
		// require.NotZero(t, transfer.ID)
		// require.NotZero(t, transfer.CreatedAt)

		// from_entry := result.FromEntry
		// require.NotEmpty(t, from_entry)
		// require.Equal(t, arg.FromAccountID,from_entry.AccountID)
		// require.Equal(t, arg.Amount,-from_entry.Amount)
		// _, err = testStore.GetEntryRecord(context.Background(), from_entry.ID)
		// require.NoError(t, err)
		// require.NotZero(t, from_entry.ID)
		// require.NotZero(t, from_entry.CreatedAt)

		// to_entry := result.ToEntry
		// require.NotEmpty(t, to_entry)
		// require.Equal(t, arg.ToAccountID,to_entry.AccountID)
		// require.Equal(t, arg.Amount,+to_entry.Amount)
		// _, err = testStore.GetEntryRecord(context.Background(), to_entry.ID)
		// require.NoError(t, err)
		// require.NotZero(t, to_entry.ID)
		// require.NotZero(t, to_entry.CreatedAt)

		// //check accounts
		// from_account := result.FromAccount
		// require.NotEmpty(t, from_account)
		// require.Equal(t, arg.FromAccountID, from_account.ID)

		// to_account := result.ToAccount
		// require.NotEmpty(t, to_account)
		// require.Equal(t, arg.ToAccountID, to_account.ID)

		// //TODO: check account's balance
		// fmt.Println(">> tx:", from_account.Balance, to_account.Balance)
		// diff1 := account1.Balance - from_account.Balance
		// diff2 := to_account.Balance - account2.Balance
		// require.Equal(t, diff1, diff2)
		// require.True(t, diff1 > 0)
		// require.True(t, diff1 % arg.Amount == 0)
		// k := diff1 / arg.Amount
		// require.True(t, k >= 1 && k <= int64(n))
		// require.NotContains(t, existed, k)
		// existed[k] = true
	}

	// check the updated final balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	// fmt.Println(">>>> after update account 1 balance", updatedAccount1.Balance)
	// fmt.Println(">>>> after update account 2 balance", updatedAccount2.Balance)
	// diff1 := account1.Balance - updatedAccount1.Balance
	// diff2 := updatedAccount2.Balance - account2.Balance
	// require.Equal(t, diff1, diff2)
	// require.True(t, diff1 > 0)
	// require.True(t, diff1 % arg.Amount == 0)
	// k := diff1 / arg.Amount
	// require.Equal(t, k, int64(n))
	// require.Equal(t, account1.Balance - arg.Amount * int64(n),updatedAccount1.Balance)
	// require.Equal(t, account2.Balance + arg.Amount * int64(n),updatedAccount2.Balance)
	require.Equal(t, account1.Balance,updatedAccount1.Balance)
	require.Equal(t, account2.Balance,updatedAccount2.Balance)
	
}