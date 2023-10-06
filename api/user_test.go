package api

import (
	"bytes"
	//"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	//"time"

	mockdb "github.com/MeganViga/BankBackend/db/mock"
	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher { return eqCreateUserParamsMatcher{arg, password} }

type eqCreateUserParamsMatcher struct {
	arg db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok{
		return false
	}
	err := util.CheckPasswordHash(e.password,arg.HashedPassword)
	if err != nil{
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
	
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}


func TestCreateUser(t *testing.T){
	user, password := randomUser(t)
	testCases := []struct{
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":user.Username,
				"password":password,
				"fullname":user.Fullname,
				"email":user.Email,
			},
			buildStubs: func(store *mockdb.MockStore){
				arg := db.CreateUserParams{
					Username: user.Username,
					Fullname: user.Fullname,
					//HashedPassword: user.HashedPassword,
					Email: user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(),EqCreateUserParams(arg,password)).
				Times(1).
				Return(user,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
				require.Equal(t,http.StatusOK,recorder.Code)
				requireBodyMatchUser(t,recorder.Body,user)
			},
		},
		// {
		// 	name: "OK",
		// 	body: gin.H{
		// 		"username":  user.Username,
		// 		"password":  password,
		// 		"full_name": user.Fullname,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.CreateUserParams{
		// 			Username: user.Username,
		// 			Fullname: user.Fullname,
		// 			Email:    user.Email,
		// 		}
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
		// 			Times(1).
		// 			Return(user, nil)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusOK, recorder.Code)
		// 		requireBodyMatchUser(t, recorder.Body, user)
		// 	},
		// },
		// {
		// 	name: "Not Found",
		// 	accountID: account.ID,
		// 	buildStubs: func(store *mockdb.MockStore){
		// 		store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).
		// 		Times(1).
		// 		Return(db.Account{},sql.ErrNoRows)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
		// 		require.Equal(t,http.StatusNotFound,recorder.Code)
		// 		//requireBodyMatchAccount(t,recorder.Body,account)
		// 	},
		// },
		// {
		// 	name: "Internel Server Error",
		// 	accountID: account.ID,
		// 	buildStubs: func(store *mockdb.MockStore){
		// 		store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account.ID)).
		// 		Times(1).
		// 		Return(db.Account{},sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
		// 		require.Equal(t,http.StatusInternalServerError,recorder.Code)
		// 		//requireBodyMatchAccount(t,recorder.Body,account)
		// 	},
		// },
		// {
		// 	name: "Bad request: Invalid ID",
		// 	accountID: 0,
		// 	buildStubs: func(store *mockdb.MockStore){
		// 		store.EXPECT().GetAccount(gomock.Any(),gomock.Any()).
		// 		Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder){
		// 		require.Equal(t,http.StatusBadRequest,recorder.Code)
		// 		// requireBodyMatchAccount(t,recorder.Body,account)
		// 	},
		// },
		// //Todo

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

	//Marshal data bidy to json
	data, err := json.Marshal(tc.body)
	require.NoError(t, err)
	url :="/users"
	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	require.NoError(t, err)
	server.router.ServeHTTP(recorder,request)

	//check response

	// require.Equal(t,http.StatusOK,recorder.Code)
	// requireBodyMatchAccount(t,recorder.Body,account)
	tc.checkResponse(t, recorder)
		})
	}
	
}

func randomUser(t *testing.T)(user db.Userdatum, password string){
	password = util.RandomName(10)
	hashedpassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.Userdatum{
		Username: util.RandomName(6),
		HashedPassword: hashedpassword,
		Fullname: util.RandomName(15),
		Email: util.RandomName(4)+ "@gmail.com",
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer,user db.Userdatum){
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var getUser db.Userdatum
	err = json.Unmarshal(data, &getUser)
	require.NoError(t, err)
	
	require.Equal(t,user.Username,getUser.Username)
	//require.Equal(t,user.HashedPassword,getUser.HashedPassword)
	require.Empty(t,getUser.HashedPassword)
	require.Equal(t,user.Fullname,getUser.Fullname)
	require.Equal(t,user.Email,getUser.Email)
	//require.Equal(t,user.PasswordChangedAt,getUser.PasswordChangedAt)
	//require.Equal(t,user.CreatedAt,getUser.CreatedAt)


}