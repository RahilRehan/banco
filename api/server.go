package api

import (
	"github.com/RahilRehan/banco/db"
	"github.com/gin-gonic/gin"
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

	router.POST("/accounts/", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers/", server.createTransfer)

	server.router = router
	return server
}
