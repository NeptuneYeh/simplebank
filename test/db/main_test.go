package db

import (
	"database/sql"
	"github.com/NeptuneYeh/simplebank/init/config"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *postgresdb.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	configModule := config.NewModule("../../")

	testDB, err = sql.Open(configModule.DBDriver, configModule.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = postgresdb.New(testDB)
	os.Exit(m.Run())
}
