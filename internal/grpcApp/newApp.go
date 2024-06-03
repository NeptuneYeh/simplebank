package grpcApp

import (
	"github.com/NeptuneYeh/simplebank/init/asynq"
	"github.com/NeptuneYeh/simplebank/init/store"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/pb"
	"github.com/NeptuneYeh/simplebank/tools/worker"
)

type Module struct {
	pb.UnimplementedSimpleBankServer
	store           postgresdb.Store
	taskDistributor worker.TaskDistributor
}

func NewModule() *Module {
	return &Module{
		store:           store.MainStore,
		taskDistributor: asynq.MainAsynq.TaskDistributor,
	}
}
