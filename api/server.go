package api

import (
	"fmt"
	"os"

	db "github.com/RahilRehan/banco/db/sqlc"
	"github.com/RahilRehan/banco/db/util"
	"github.com/RahilRehan/banco/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
	symmetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")
	tokenMaker, err := token.NewPasetoMaker(symmetricKey)
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
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts/", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts/", server.listAccounts)
	authRoutes.PUT("/accounts/", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers/", server.createTransfer)

	router.POST("/users/", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
}
