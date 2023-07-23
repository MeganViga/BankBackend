package db

import (
	"context"
	"testing"

	"github.com/MeganViga/BankBackend/util"
	"github.com/stretchr/testify/require"
)
func createRandomTransferRecord(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTranferRecordParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: int64(util.RandomNumber(1, int(account1.Balance))),
	}
	transfer, err := testQueries.CreateTranferRecord(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t,transfer.CreatedAt)

}
func createRandomTransferRecordByFromAndToAccount(t *testing.T, account1, account2 Account){
	arg := CreateTranferRecordParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: int64(util.RandomNumber(1, int(account1.Balance))),
	}
	transfer, err := testQueries.CreateTranferRecord(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t,transfer.CreatedAt)

}
func TestGetTransferRecord(t *testing.T){
	createRandomTransferRecord(t)
}

func TestListTransferRecords(t *testing.T){
	for i :=0 ; i< 10; i++{
		createRandomTransferRecord(t)
	}
	arg := ListTransferRecordsParams{
		Limit: 5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransferRecords(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers{
		require.NotEmpty(t, transfer)
	}
}

func TestListTransferRecordsByFromAndToAccountID(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	for i :=0 ; i< 10; i++{
		createRandomTransferRecordByFromAndToAccount(t, account1, account2)
	}
	arg := ListTransferRecordsByFromAndToAccountIDParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Limit: 5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransferRecordsByFromAndToAccountID(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers{
		require.NotEmpty(t, transfer)
	}
	}
