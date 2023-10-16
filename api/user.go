package api

import (
	"database/sql"
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
type userResponse struct{
	Username          string    `json:"username"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
func newUserResponse(user db.Userdatum)userResponse{
	return userResponse{
		Username: user.Username,
		Fullname: user.Fullname,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
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
	//arg = db.CreateUserParams{}
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
	rsp := newUserResponse(user)

	ctx.JSON(http.StatusOK,rsp)
}


type loginUserRequest struct{
	Username   string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}
type loginUserResponse struct{
	AccessToken  string `json:"access_token"`
	User userResponse	`json:"user"`
}

func (s *Server)loginUser(ctx *gin.Context){
	var r loginUserRequest
	if err := ctx.ShouldBindJSON(&r); err != nil{
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return
	}
	userdata, err := s.store.GetUser(ctx, r.Username)
	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}
	if err := util.CheckPasswordHash(r.Password,userdata.HashedPassword); err != nil{
		ctx.JSON(http.StatusUnauthorized,errResponse(err))
		return
	}
	access_token, err := s.tokenMaker.CreateToken(r.Username,s.config.AccessTokenDuration)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errResponse(err))
		return
	}
	
	rsp := loginUserResponse{
		AccessToken: access_token,
		User: newUserResponse(userdata),
	}
	ctx.JSON(http.StatusOK,rsp)
	
}