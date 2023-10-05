package api

import (
	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	//"github.com/vektah/gqlparser/v2/validator/v10"
)

//Server serves requests for http bank service
type Server struct{
	store db.Store
	router *gin.Engine
}


func NewServer(store db.Store)*Server{
	server := &Server{
		store: store,
	}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate);ok{
			v.RegisterValidation("currency",validCurrency)
			// v.RegisterValidation("email",validEmail) --> gin validator has inbuil email validator, that's why commenting this
	}
	router.POST("/accounts",server.createAccount)
	router.GET("/account/:id",server.getAccount)
	router.GET("/accounts",server.getAccounts)
	router.DELETE("/account/:id",server.deleteAccount)
	router.POST("/transfers",server.createTransfer)
	router.POST("/users",server.createUser)

	server.router = router
	return server
}

func(s *Server)StartServer(address string)error{
	return s.router.Run(address)
}

func errResponse(err error)gin.H{
	return gin.H{"error":err.Error()}
}