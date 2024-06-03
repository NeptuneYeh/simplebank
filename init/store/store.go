package store

import (
	"database/sql"
	"errors"
	"github.com/NeptuneYeh/simplebank/init/config"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// TODO run db migration
	runDBMigration(config.MainConfig.MigrationURL, config.MainConfig.DBSource)

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

func runDBMigration(migrationURL string, dbSource string) {
	// TODO run db migration
	// TODO 應該要實作 prod phase 以外才會真正執行, 不然很危險
	if config.MainConfig.ENV != "prod" {
		migration, err := migrate.New(migrationURL, dbSource)
		if err != nil {
			log.Fatal("cannot create migration: ", err)
		}

		if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("cannot run migration up: ", err)
		}

		log.Println("migration completed")
	}
}
