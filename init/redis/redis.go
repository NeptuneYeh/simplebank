package redis

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/redis/go-redis/v9"
	"log"
)

var MyRedisClient *redis.Client

type Module struct {
	Client *redis.Client
}

func NewModule() *Module {
	client := redis.NewClient(&redis.Options{
		Addr:     config.MainConfig.RedisAddress,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("cannot connect to redis: ", err)
	}

	MyRedisClient = client
	redisModule := &Module{
		Client: client,
	}

	return redisModule
}
