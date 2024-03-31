package api

import (
	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	docs "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/docs"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {

	docs.SwaggerInfo.BasePath = "/api/v1"

	server := &Server{store: store}
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	

	v1 := router.Group("/api/v1")		
	{
		v1.POST("/accounts", server.CreateAccount)
		v1.GET("/accounts/:id", server.GetAccount)
		v1.GET("/accounts", server.GetAccounts)
		v1.DELETE("/accounts/:id", server.DeleteAccount)
		v1.POST("/accounts/deposit", server.Deposit)

		//transfer
		v1.POST("/transfers", server.CreateTransfer)
		v1.GET("/transfers", server.GetTransfersByAccount)
		v1.GET("/transfers/:id", server.GetTransferById)

		//entry
		v1.GET("/entry", server.GetEntriesByAccount)

		//user
		v1.POST("/users/register", server.CreateUser)
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
