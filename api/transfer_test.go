package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/MeganViga/BankBackend/db/mock"
	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T){
	user1, _ := randomUser(t)
	account1 := randomAccount(user1.Username)
	account1.Balance = 200
	account1.Currency = util.USD
	user2, _ := randomUser(t)
	account2 := randomAccount(user2.Username)
	account2.Balance = 100
	account2.Currency = util.USD
	transactionAmount := int64(10)
	//result := transferResult(account1, account2,transactionAmount)
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)
	//build stubs
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

	arg := db.TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        transactionAmount,
			}
	store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
	//start the http test server
	server := newTestServer(t,store)
	recorder := httptest.NewRecorder()
	url :="/transfers"
	arg2 := transferRequest{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: transactionAmount,
		Currency: util.USD,
	}
	
	bytess, err := json.Marshal(arg2)
	require.NoError(t, err)
	body:= bytes.NewReader(bytess)
	request, err := http.NewRequest("POST", url, body)
	require.NoError(t, err)
	addAuthorization(t, request, server.tokenMaker,authorizationTypeBearer,user1.Username, time.Minute )
	server.router.ServeHTTP(recorder,request)

	//check response
	require.Equal(t,http.StatusOK,recorder.Code)
	//requireBodyMatchAccount(t, recorder.Body, result)
}

// func transferResult(account1, account2 db.Account, amount int64)db.TransferTxResult{
// 	account1.Balance = account1.Balance - amount
// 	account2.Balance = account2.Balance + amount
// 	return db.TransferTxResult{
// 		Transfer: db.Transfer{
// 			ID: int64(util.RandomNumber(1, 1000)),
// 			FromAccountID: account1.ID,
// 			ToAccountID: account2.ID,
// 			Amount: amount,
// 			CreatedAt: time.Now(),
// 		},
// 		FromAccount: account1,
// 		ToAccount: account2,
// 		FromEntry: db.Entry{
// 			ID: int64(util.RandomNumber(1, 1000)),
// 			AccountID: account1.ID,
// 			Amount: -amount,
// 			CreatedAt: time.Now(),
// 		},
// 		ToEntry: db.Entry{
// 			ID: int64(util.RandomNumber(1, 1000)),
// 			AccountID: account2.ID,
// 			Amount: amount,
// 			CreatedAt: time.Now(),
// 		},
		
// 	}
// }

// func randomAccountWithCurrency(owner string, currency string)db.Account{
// 	return db.Account{
// 		ID: int64(util.RandomNumber(1, 1000)),
// 		Owner: owner,
// 		Balance: int64(util.RandomNumber(10,100)),
// 		Currency: currency,
// 		CreatedAt: time.Now(),
// 	}
// }