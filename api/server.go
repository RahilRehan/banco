package api

import (
	"fmt"

	db "github.com/RahilRehan/banco/db/sqlc"
	"github.com/RahilRehan/banco/db/util"
	"github.com/RahilRehan/banco/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func (s *server) Start(address string) error {
	return s.router.Run(address)
}

func NewServer(cfg util.Config, store db.Store) (*server, error) {
	tokenMaker, err := token.NewPasetoMaker(viper.GetString("TOKEN_SYMMETRIC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     cfg,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *server) setupRouter() {
	router := gin.Default()
	router.POST("/accounts/", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccounts)
	router.PUT("/accounts/", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers/", server.createTransfer)

	router.POST("/users/", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}
