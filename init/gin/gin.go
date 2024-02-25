package gin

import (
	"github.com/NeptuneYeh/simplebank/internal/application/controllers"
	"github.com/gin-gonic/gin"
)

type Module struct {
	router *gin.Engine
}

func NewModule() *Module {
	r := gin.Default()
	ginModule := &Module{
		router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute()

	return ginModule
}

// setup route
func (module *Module) setupRoute() {
	// init controller
	accountController := controllers.NewAccountController()
	// add routes to router
	module.router.POST("/accounts", accountController.CreateAccount)
	module.router.GET("/accounts/:id", accountController.GetAccount)
	module.router.GET("/accounts", accountController.ListAccount)
}

// Run gin
func (module *Module) Run(address string) {
	err := module.router.Run(address)
	if err != nil {
		return
	}
}
