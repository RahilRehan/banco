package api

import (
	db "github.com/RahilRehan/banco/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type server struct {
	store  db.Store
	router *gin.Engine
}

func (s *server) Start(address string) error {
	return s.router.Run(address)
}

func NewServer(store db.Store) *server {
	server := &server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.POST("/accounts/", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers/", server.createTransfer)

	server.router = router
	return server
}
