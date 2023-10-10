package token

import (
	"testing"
	"time"

	"github.com/MeganViga/BankBackend/util"
	"github.com/stretchr/testify/require"
)


func TestPasetoMaker(t *testing.T){
	maker, err :=NewPasetoMaker(util.RandomName(32))
	require.NoError(t, err)
	username := util.RandomName(6)
	duration := time.Minute * 1
	issued_at := time.Now()
	expired_at := issued_at.Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issued_at, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expired_at, payload.ExpiredAt, time.Second)
	
}

func TestExpiredPasetoToken(t *testing.T){
	maker, err := NewPasetoMaker(util.RandomName(32))
	require.NoError(t, err)
	token, err := maker.CreateToken(util.RandomName(6), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t, err, ErrTokenExpired.Error())
	require.Nil(t, payload)

}

//None algo doesn't exist in paseto