package configs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Rc *redis.Client

func InitRedis() {
	Rc = redis.NewClient(&redis.Options{
		Addr: Conf.Redis.Addr,
		DB:   Conf.Redis.Db,
	})
	_, err := Rc.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis init success")
}
