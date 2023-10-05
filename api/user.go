package api

import (
	"net/http"
	"time"

	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct{
	Username   string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	Fullname string `json:"fullname" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
type createUserResponse struct{
	Username          string    `json:"username"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
func (s *Server)createUser(ctx *gin.Context){
	var r createUserRequest
	if err := ctx.ShouldBindJSON(&r); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(r.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username: r.Username,
		HashedPassword: hashedPassword,
		Fullname: r.Fullname,
		Email: r.Email,

	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil{
		if pqErr := err.(*pq.Error);pqErr != nil{
			errName := pqErr.Code.Name()
			switch errName{
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	rsp := createUserResponse{
		Username: user.Username,
		Fullname: user.Fullname,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK,rsp)
}