package config

import (
	"github.com/spf13/viper"
	"log"
)

type Module struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func NewModule() *Module {
	// 這裡的 "./" 代表目前工作目錄，也就是你在命令列中執行 Go 指令時所處的目錄。
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var configModule Module
	err = viper.Unmarshal(&configModule)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &configModule
}
