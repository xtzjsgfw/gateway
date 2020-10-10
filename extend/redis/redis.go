package redis

import (
	"fmt"
	"gateway/extend/conf"
	"github.com/go-redis/redis"
)

var Rdb *redis.Client
var err error

func Init() {
	addr := fmt.Sprintf("%s:%s", conf.RedisConf.Host, conf.RedisConf.Port)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
	})
	pone, err := Rdb.Ping().Result()
	if err != nil {
		fmt.Println("连接redis失败")
	}
	fmt.Println("连接成功", pone)
}