package auth

import (
	"github.com/NeptuneYeh/simplebank/init/config"
	myToken "github.com/NeptuneYeh/simplebank/tools/token"
	"log"
)

var MainAuth myToken.Maker

type Module struct {
	TokenMaker myToken.Maker
}

func NewModule() *Module {
	maker, err := myToken.NewPasetoMaker(config.MainConfig.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("Failed to create token maker: %v", err)
	}

	MainAuth = maker
	authModule := &Module{
		TokenMaker: maker,
	}

	return authModule
}
