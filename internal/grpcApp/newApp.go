package grpcApp

import (
	"github.com/NeptuneYeh/simplebank/init/store"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/pb"
)

type Module struct {
	pb.UnimplementedSimpleBankServer
	store postgresdb.Store
}

func NewModule() *Module {
	return &Module{
		store: store.MainStore,
	}
}
