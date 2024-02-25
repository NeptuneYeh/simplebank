package gin

import (
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/internal/application/controllers"
	"github.com/gin-gonic/gin"
	"log"
)

type Module struct {
	router *gin.Engine
}

func NewModule(config *config.Module) *Module {
	r := gin.Default()
	ginModule := &Module{
		router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute(config)

	return ginModule
}

// setup route
func (module *Module) setupRoute(config *config.Module) {
	// init controller
	accountController := controllers.NewAccountController(config)
	// add routes to router
	module.router.POST("/accounts", accountController.CreateAccount)
	module.router.GET("/accounts/:id", accountController.GetAccount)
	module.router.GET("/accounts", accountController.ListAccount)
}

// Run gin
func (module *Module) Run(address string) {
	err := module.router.Run(address)
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
