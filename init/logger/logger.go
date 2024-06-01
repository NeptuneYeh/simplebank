package logger

import (
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var MainLog *zerolog.Logger

type Module struct {
	Logger zerolog.Logger
}

func NewModule() *Module {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if config.MainConfig.ENV != "prod" {
		// global pretty logging 效果
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// 如果只想要個別創建的實例有 pretty logging 效果
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFormatUnix}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	}

	MainLog = &logger

	return &Module{
		Logger: logger,
	}
}
