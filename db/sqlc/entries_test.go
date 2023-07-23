package db

import (
	"context"
	"testing"

	"github.com/MeganViga/BankBackend/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntryRecord(t *testing.T){
	account1 := createRandomAccount(t)
	arg :=  CreateEntryRecordParams{
		AccountID: account1.ID,
		Amount: int64(util.RandomNumber(-int(account1.Balance), int(account1.Balance))),
	}
	entry , err := testQueries.CreateEntryRecord(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t,entry.CreatedAt)
}

func createRandomEntryRecordForAccount(t *testing.T, account Account){
	arg :=  CreateEntryRecordParams{
		AccountID: account.ID,
		Amount: int64(util.RandomNumber(-int(account.Balance), int(account.Balance))),
	}
	entry , err := testQueries.CreateEntryRecord(context.Background(),arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t,entry.CreatedAt)
}
func TestGetEntryRecord(t *testing.T){
	createRandomEntryRecord(t)
}

func TestListEntryRecords(t *testing.T){
	for i :=0 ; i< 10; i++{
		createRandomEntryRecord(t)
	}
	arg := ListEntryRecordsParams{
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntryRecords(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries{
		require.NotEmpty(t, entry)
	}
}


func TestListEntryRecordsByAccountID(t *testing.T){
	account := createRandomAccount(t)
	for i :=0 ; i< 10; i++{
		createRandomEntryRecordForAccount(t, account)
	}
	arg := ListEntryRecordsByAccountIDParams{
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntryRecordsByAccountID(context.Background(),arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries{
		require.NotEmpty(t, entry)
	}
}