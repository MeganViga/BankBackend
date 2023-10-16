package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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


func TestGetAccount(t *testing.T){
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct{
		name string
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(account,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusOK,recorder.Code)
				requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "Not Found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{},sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusNotFound,recorder.Code)
				//requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "Internel Server Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
				//requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "Bad request: Invalid ID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore){
				store.EXPECT().GetAccount(gomock.Any(),gomock.Any()).
				Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusBadRequest,recorder.Code)
				// requireBodyMatchAccount(t,recorder.Body,account)
			},
		},	
		//Todo

	}
	for i := range testCases{
		tc := testCases[i]
		t.Run(tc.name,func(t *testing.T){
			ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)
	//build stubs
	// store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).
	// Times(1).
	// Return(account,nil)
	tc.buildStubs(store)

	//start test http server

	server := newTestServer(t,store)
	recorder := httptest.NewRecorder()
	url :=fmt.Sprintf("/account/%d", tc.accountID)
	request, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	addAuthorization(t, request, server.tokenMaker,authorizationTypeBearer, user.Username, time.Minute)
	server.router.ServeHTTP(recorder,request)

	//check response

	// require.Equal(t,http.StatusOK,recorder.Code)
	// requireBodyMatchAccount(t,recorder.Body,account)
	tc.checkResponse(t, recorder)
		})
	}
	
}

func TestCreateAccount(t *testing.T){
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	account.Balance = 0
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	//build stubs
	arg := db.CreateAccountParams{
		Owner: account.Owner,
		Balance: 0,
		Currency: account.Currency,
	}

	store.EXPECT().CreateAccount(gomock.Any(),gomock.Eq(arg)).Times(1).Return(account,nil)
	//start the http test server
	server := newTestServer(t,store)
	recorder := httptest.NewRecorder()
	url :="/accounts"
	arg2 := createAccountRequest{
		//Owner: arg.Owner,
		Currency: arg.Currency,
	}
	bytess, err := json.Marshal(arg2)
	require.NoError(t, err)
	body:= bytes.NewReader(bytess)
	request, err := http.NewRequest("POST", url, body)
	require.NoError(t, err)
	addAuthorization(t, request, server.tokenMaker,authorizationTypeBearer,user.Username, time.Minute )
	server.router.ServeHTTP(recorder,request)

	//check response
	require.Equal(t,http.StatusOK,recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)


}

func randomAccount(owner string)db.Account{
	return db.Account{
		ID: int64(util.RandomNumber(1, 1000)),
		Owner: owner,
		Balance: int64(util.RandomNumber(10,100)),
		Currency: util.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer,account db.Account){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	
	require.Equal(t,account.ID,getAccount.ID)
	require.Equal(t,account.Owner,getAccount.Owner)
	require.Equal(t,account.Balance,getAccount.Balance)
	require.Less(t,account.CreatedAt,time.Now())


}