package asynq

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/config"
	db "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/mail"
	"github.com/NeptuneYeh/simplebank/tools/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"time"
)

var MainAsynq *Module

type Module struct {
	RedisOpt        *asynq.RedisClientOpt
	TaskDistributor worker.TaskDistributor
	TaskProcessor   worker.TaskProcessor
}

func NewModuleForTest(mockTaskDistributor worker.TaskDistributor) *Module {
	redisOpt := &asynq.RedisClientOpt{
		Addr: config.MainConfig.RedisAddress,
	}

	asynqModule := &Module{
		TaskDistributor: mockTaskDistributor,
		RedisOpt:        redisOpt,
	}
	MainAsynq = asynqModule

	return asynqModule
}

func NewModule() *Module {
	redisOpt := &asynq.RedisClientOpt{
		Addr: config.MainConfig.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(*redisOpt)

	asynqModule := &Module{
		TaskDistributor: taskDistributor,
		RedisOpt:        redisOpt,
	}
	MainAsynq = asynqModule

	return asynqModule
}

func (module *Module) Run(redisOpt *asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.MainConfig.EmailSenderName, config.MainConfig.EmailSenderAddress, config.MainConfig.EmailSenderPassword)
	module.TaskProcessor = worker.NewRedisTaskProcessor(*redisOpt, store, mailer)
	go func() {
		err := module.TaskProcessor.Start()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to start task processor")
		}
	}()
}

func (module *Module) Shutdown() error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	module.TaskProcessor.Shutdown()
	return nil
}
