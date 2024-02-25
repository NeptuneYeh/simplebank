package gin

import (
	"github.com/NeptuneYeh/simplebank/internal/application/controllers"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

type Module struct {
	Router *gin.Engine
}

func NewModule(store postgresdb.Store) *Module {
	r := gin.Default()
	ginModule := &Module{
		Router: r,
	}
	gin.ForceConsoleColor()
	ginModule.setupRoute(store)

	return ginModule
}

// setup route
func (module *Module) setupRoute(store postgresdb.Store) {
	// init controller
	accountController := controllers.NewAccountController(store)
	// add routes to router
	module.Router.POST("/accounts", accountController.CreateAccount)
	module.Router.GET("/accounts/:id", accountController.GetAccount)
	module.Router.GET("/accounts", accountController.ListAccount)
}

// Run gin
func (module *Module) Run(address string) {
	err := module.Router.Run(address)
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
