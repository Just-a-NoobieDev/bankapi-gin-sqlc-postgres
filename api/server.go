package api

import (
	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	

	v1 := router.Group("/api/v1")		
	{
		v1.POST("/accounts", server.CreateAccount)
		v1.GET("/accounts/:id", server.GetAccount)
		v1.GET("/accounts", server.GetAccounts)
		v1.DELETE("/accounts/:id", server.DeleteAccount)
		v1.POST("/accounts/deposit", server.Deposit)
	}

	

	server.router = router
	
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}


func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
