package init

import (
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/gapi"
	"github.com/NeptuneYeh/simplebank/init/grpcGateway"
	"github.com/NeptuneYeh/simplebank/init/logger"
	"github.com/NeptuneYeh/simplebank/init/redis"
	"github.com/NeptuneYeh/simplebank/init/store"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"syscall"
)

type MainInitProcess struct {
	ConfigModule *config.Module
	LogModule    *logger.Module
	AuthModule   *auth.Module
	StoreModule  *store.Module
	RedisModule  *redis.Module
	//GinModule         *gin.Module
	GRPCModule        *gapi.Module
	GRPCGatewayModule *grpcGateway.Module
	OsChannel         chan os.Signal
}

func NewMainInitProcess(configPath string) *MainInitProcess {
	configModule := config.NewModule(configPath)
	logModule := logger.NewModule()
	authModule := auth.NewModule()
	storeModule := store.NewModule()
	redisModule := redis.NewModule()
	//ginModule := gin.NewModule()
	gapiModule := gapi.NewModule()
	grpcGatewayModule := grpcGateway.NewModule()

	channel := make(chan os.Signal, 1)
	return &MainInitProcess{
		ConfigModule: configModule,
		LogModule:    logModule,
		AuthModule:   authModule,
		StoreModule:  storeModule,
		RedisModule:  redisModule,
		//GinModule:         ginModule,
		GRPCModule:        gapiModule,
		GRPCGatewayModule: grpcGatewayModule,
		OsChannel:         channel,
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	//m.GinModule.Run(m.ConfigModule.ServerAddress)
	m.GRPCGatewayModule.Run(m.ConfigModule.ServerAddress)
	m.GRPCModule.Run(m.ConfigModule.GRPCServerAddress)
	// register os signal for graceful shutdown
	signal.Notify(m.OsChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-m.OsChannel
	m.LogModule.Logger.Fatal().Msg("Received signal: " + s.String())
	//_ = m.GinModule.Shutdown()
	_ = m.GRPCGatewayModule.Shutdown()
	_ = m.GRPCModule.Shutdown()
}
