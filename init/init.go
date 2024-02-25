package init

import ginModule "github.com/NeptuneYeh/simplebank/init/gin"

type MainInitProcess struct {
	ginModule *ginModule.Module
}

func NewMainInitProcess() *MainInitProcess {
	return &MainInitProcess{
		ginModule: ginModule.NewModule(),
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	m.ginModule.Run("0.0.0.0:8080")
}
