package init

import (
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/gin"
)

type MainInitProcess struct {
	configModule *config.Module
	ginModule    *gin.Module
}

func NewMainInitProcess() *MainInitProcess {
	configModule := config.NewModule()
	return &MainInitProcess{
		configModule: configModule,
		ginModule:    gin.NewModule(configModule),
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	m.ginModule.Run(m.configModule.ServerAddress)
}
