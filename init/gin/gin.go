package gin

import (
	"context"
	"github.com/NeptuneYeh/simplebank/internal/application/controllers"
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

func NewModule() *Module {
	r := gin.Default()
	ginModule := &Module{
		Router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute()

	return ginModule
}

// setup route
func (module *Module) setupRoute() {
	// init controller
	accountController := controllers.NewAccountController()
	transferController := controllers.NewTransferController()
	// binding validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", controllers.ValidCurrency)
		if err != nil {
			return
		}
	}
	// add routes to router
	module.Router.POST("/accounts", accountController.CreateAccount)
	module.Router.GET("/accounts/:id", accountController.GetAccount)
	module.Router.GET("/accounts", accountController.ListAccount)

	module.Router.POST("/transfers", transferController.CreateTransfer)
}

// Run gin
func (module *Module) Run(address string) {
	module.Server = &http.Server{
		Addr:    address,
		Handler: module.Router,
	}

	go func() {
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
