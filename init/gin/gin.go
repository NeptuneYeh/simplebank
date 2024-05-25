package gin

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/internal/application/controllers"
	"github.com/NeptuneYeh/simplebank/internal/application/middlewares"
	myValidator "github.com/NeptuneYeh/simplebank/tools/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

type Module struct {
	Router *gin.Engine
	Server *http.Server
}

var MainGin *Module

func NewModule() *Module {
	r := gin.Default()
	ginModule := &Module{
		Router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute()

	MainGin = ginModule

	return ginModule
}

// setup route
func (module *Module) setupRoute() {
	// init controller
	userController := controllers.NewUserController()
	accountController := controllers.NewAccountController()
	transferController := controllers.NewTransferController()
	// binding validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", myValidator.ValidCurrency)
		if err != nil {
			return
		}
	}

	// add routes to router
	module.Router.POST("/users", userController.CreateUser)
	module.Router.POST("/users/login", userController.LoginUser)
	module.Router.POST("/tokens/renew_access", userController.RenewAccessToken)

	authRoutes := module.Router.Group("/").Use(middlewares.AuthMiddleware(auth.MainAuth))
	authRoutes.POST("/accounts", accountController.CreateAccount)
	authRoutes.GET("/accounts/:id", accountController.GetAccount)
	authRoutes.GET("/accounts", accountController.ListAccount)

	authRoutes.POST("/transfers", transferController.CreateTransfer)
}

// Run gin
func (module *Module) Run(address string) {
	module.Server = &http.Server{
		Addr:    address,
		Handler: module.Router,
	}

	go func() {
		log.Printf("Starting gin framework server on %s\n", address)
		if err := module.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
}

func (module *Module) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := module.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to run Gin shutdown: %v", err)
	}
	return nil
}
