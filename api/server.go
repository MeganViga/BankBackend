package api

import (
	db "github.com/MeganViga/BankBackend/db/sqlc"
	"github.com/gin-gonic/gin"
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
	router.POST("/accounts",server.createAccount)
	router.GET("/account/:id",server.getAccount)
	router.GET("/accounts",server.getAccounts)
	router.DELETE("/account/:id",server.deleteAccount)

	server.router = router
	return server
}

func(s *Server)StartServer(address string)error{
	return s.router.Run(address)
}

func errResponse(err error)gin.H{
	return gin.H{"error":err.Error()}
}