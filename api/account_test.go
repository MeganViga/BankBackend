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

	mockdb "github.com/MeganViga/BankBackend/db/mock"
	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)


func TestGetAccount(t *testing.T){
	account := randomAccount()
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

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url :=fmt.Sprintf("/account/%d", tc.accountID)
	request, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder,request)

	//check response

	// require.Equal(t,http.StatusOK,recorder.Code)
	// requireBodyMatchAccount(t,recorder.Body,account)
	tc.checkResponse(t, recorder)
		})
	}
	
}

func randomAccount()db.Account{
	return db.Account{
		ID: int64(util.RandomNumber(1, 1000)),
		Owner: util.RandomName(7),
		Balance: int64(util.RandomNumber(10,100)),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer,account db.Account){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	require.Equal(t,account,getAccount)

}