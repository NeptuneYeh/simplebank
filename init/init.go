package init

import (
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/gin"
	"github.com/NeptuneYeh/simplebank/init/logger"
	"github.com/NeptuneYeh/simplebank/init/store"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

type MainInitProcess struct {
	ConfigModule *config.Module
	LogModule    *logger.Module
	StoreModule  *store.Module
	GinModule    *gin.Module
	OsChannel    chan os.Signal
}

func NewMainInitProcess(configPath string) *MainInitProcess {
	configModule := config.NewModule(configPath)
	logModule := logger.NewModule()
	storeModule := store.NewModule()
	ginModule := gin.NewModule()

	channel := make(chan os.Signal, 1)
	return &MainInitProcess{
		ConfigModule: configModule,
		LogModule:    logModule,
		StoreModule:  storeModule,
		GinModule:    ginModule,
		OsChannel:    channel,
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	m.GinModule.Run(m.ConfigModule.ServerAddress)
	// register os signal for graceful shutdown
	signal.Notify(m.OsChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-m.OsChannel
	m.LogModule.Logger.Fatal().Msg("Received signal: " + s.String())
	_ = m.GinModule.Shutdown()
}
