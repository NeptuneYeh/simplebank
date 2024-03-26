package store

import (
	"database/sql"
	"github.com/NeptuneYeh/simplebank/init/config"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	_ "github.com/lib/pq"
	"log"
)

var MainStore postgresdb.Store

type Module struct {
	Store postgresdb.Store
}

func NewModule() *Module {
	// init Store
	conn, err := sql.Open(config.MainConfig.DBDriver, config.MainConfig.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := postgresdb.NewStore(conn)
	MainStore = store

	storeModule := &Module{
		Store: store,
	}

	return storeModule
}

func NewModuleForTest(store postgresdb.Store) *Module {
	// init Store
	MainStore = store
	storeModule := &Module{
		Store: store,
	}

	return storeModule
}
