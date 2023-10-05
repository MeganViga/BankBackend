package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/MeganViga/BankBackend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T)Userdatum{
	hashedPassword, err := util.HashPassword(util.RandomName(15))
	require.NoError(t, err)
	c := CreateUserParams{
		Username: util.RandomName(6),
		HashedPassword: hashedPassword,
		Fullname: util.RandomName(10),
		Email: fmt.Sprintf("%s@gmail.com", util.RandomName(5)),
	}
	user, err := testQueries.CreateUser(context.Background(),c )

	require.NoError(t, err)
	require.NotEmpty(t,user)
	require.Equal(t,user.Username, c.Username)
	require.Equal(t,user.HashedPassword, c.HashedPassword)
	require.Equal(t,user.Fullname, c.Fullname)
	require.Equal(t,user.Email, c.Email)

	require.True(t,user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user

}


func TestCreateUser(t *testing.T){
	createRandomUser(t)
}

func TestGetUser(t *testing.T){
	user1 := createRandomUser(t)
	user2 , err := testQueries.GetUser(context.Background(),user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	
	require.Equal(t,user1.Username, user2.Username)
	require.Equal(t,user1.HashedPassword, user2.HashedPassword)
	require.Equal(t,user1.Fullname, user2.Fullname)
	require.Equal(t,user1.Email, user2.Email)
	require.True(t,user1.PasswordChangedAt.IsZero())
	require.NotZero(t, user1.CreatedAt)
	require.WithinDuration(t,user1.PasswordChangedAt, user2.PasswordChangedAt,time.Second)
	require.WithinDuration(t,user1.CreatedAt, user2.CreatedAt,time.Second)

	

}