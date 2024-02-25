package init

import (
	"database/sql"
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/gin"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	_ "github.com/lib/pq"
	"log"
)

type MainInitProcess struct {
	configModule *config.Module
	storeModule  postgresdb.Store
	ginModule    *gin.Module
}

func NewMainInitProcess() *MainInitProcess {
	configModule := config.NewModule()
	// init Store
	conn, err := sql.Open(configModule.DBDriver, configModule.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := postgresdb.NewStore(conn)
	return &MainInitProcess{
		configModule: configModule,
		storeModule:  store,
		ginModule:    gin.NewModule(store),
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	m.ginModule.Run(m.configModule.ServerAddress)
}
