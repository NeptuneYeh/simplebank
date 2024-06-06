package db

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/config"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testStore postgresdb.Store

func TestMain(m *testing.M) {
	var err error
	configModule := config.NewModule("../../")

	connPool, err := pgxpool.New(context.Background(), configModule.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	//testStore = postgresdb.New(connPool)
	testStore = postgresdb.NewStore(connPool)
	os.Exit(m.Run())
}
