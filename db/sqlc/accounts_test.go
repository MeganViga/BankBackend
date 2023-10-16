package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/MeganViga/BankBackend/util"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T)Account{
	u := createRandomUser(t)
	c := CreateAccountParams{
		Owner: u.Username,
		Balance: int64(util.RandomNumber(100, 1000)),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), c)

	require.NoError(t, err)
	require.NotEmpty(t,account)
	require.Equal(t,account.Owner, c.Owner)
	require.Equal(t,account.Balance, c.Balance)
	require.Equal(t,account.Currency, c.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account

}


func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account1 := createRandomAccount(t)
	account2 , err := testQueries.GetAccount(context.Background(),account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotZero(t, account2.ID)
	require.NotZero(t, account2.CreatedAt)
	require.Equal(t,account1.Owner, account2.Owner)
	require.Equal(t,account1.Balance, account2.Balance)
	require.Equal(t,account1.Currency, account2.Currency)

}

func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID: account1.ID,
		Balance: int64(util.RandomNumber(100, 1000)),
	}
	_, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	account2 , err :=  testQueries.GetAccount(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.NotZero(t, account2.CreatedAt)
	//t.Log(account1, account2)
}

func TestDeleteAccount(t *testing.T){
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2 , err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	// assert.Equal(t, "sql: no rows in result set", err.Error())
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T){
	var lastAccount Account
	for i :=0 ; i< 10; i++{
		lastAccount = createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	//require.Len(t, accounts, 5)
	for _, account := range accounts{
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}