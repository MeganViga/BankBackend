package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency" `
}
func (s *Server)createAccount(ctx *gin.Context){
	var r createAccountRequest
	if err := ctx.ShouldBindJSON(&r); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner: r.Owner,
		Currency: r.Currency,
		Balance: 0,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil{
		if pqErr := err.(*pq.Error);pqErr != nil{
			errName := pqErr.Code.Name()
			switch errName{
			case "foreign_key_violation","unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}
type getAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}
func (s *Server)getAccount(ctx *gin.Context){
	var r getAccountRequest
	if err := ctx.ShouldBindUri(&r); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	fmt.Println(r.ID)
	account, err := s.store.GetAccount(ctx,r.ID)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// account = db.Account{}

	ctx.JSON(http.StatusOK,account)
}


type getAccountsRequest struct{
	PageID int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}
func (s *Server)getAccounts(ctx *gin.Context){
	var r getAccountsRequest
	if err := ctx.ShouldBindQuery(&r); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	arg := db.ListAccountsParams{
		Limit: int32(r.PageSize),
		Offset: int32(r.PageID - 1) * int32(r.PageSize),
	}
	
	accounts, err := s.store.ListAccounts(ctx,arg)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,accounts)
}
type deleteAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}
func (s *Server)deleteAccount(ctx *gin.Context){
	var r deleteAccountRequest
	if err := ctx.ShouldBindUri(&r); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	fmt.Println(r.ID)
	err := s.store.DeleteAccount(ctx,r.ID)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"message":"Deleted Successfully"})
}
