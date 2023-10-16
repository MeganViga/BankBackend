package api

import (
	"fmt"

	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/MeganViga/BankBackend/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	//"github.com/vektah/gqlparser/v2/validator/v10"
	"github.com/MeganViga/BankBackend/token"
)

//Server serves requests for http bank service
type Server struct{
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}


func NewServer(config util.Config,store db.Store)(*Server, error){
	tokenMaker, err:= token.NewPasetoMaker(config.TokenSymmetricKey)
	// tokenMaker, err:= token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil{
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate);ok{
			v.RegisterValidation("currency",validCurrency)
			// v.RegisterValidation("email",validEmail) --> gin validator has inbuil email validator, that's why commenting this
	}
	server.setupRouter()
	// router.POST("/accounts",server.createAccount)
	// router.GET("/account/:id",server.getAccount)
	// router.GET("/accounts",server.getAccounts)
	// router.DELETE("/account/:id",server.deleteAccount)
	// router.POST("/transfers",server.createTransfer)
	// router.POST("/users",server.createUser)
	// router.POST("/login",server.loginUser)

	// server.router = router
	return server, nil
}

func (s *Server)setupRouter(){

	router := gin.Default()
	router.POST("/users",s.createUser)
	router.POST("/users/login",s.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts",s.createAccount)
	authRoutes.GET("/account/:id",s.getAccount)
	authRoutes.GET("/accounts",s.getAccounts)
	authRoutes.DELETE("/account/:id",s.deleteAccount)
	authRoutes .POST("/transfers",s.createTransfer)
	s.router = router
}

func(s *Server)StartServer(address string)error{
	return s.router.Run(address)
}

func errResponse(err error)gin.H{
	return gin.H{"error":err.Error()}
}