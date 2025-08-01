package config

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/kuqsi/exchangeapp/global"
)

func InitRedis() {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "192.168.189.129:6380",
		DB:       0,
		Password: "123456",
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("reids连接失败,error is %v", err)
	}

	global.RedisDb = RedisClient
}
