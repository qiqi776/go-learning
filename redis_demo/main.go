package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

// 初始化 Redis 连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
	PoolSize: 20,
	})

	_, err = rdb.Ping().Result()
	return err
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed, err:%v\n", err)
		return 
	}
	fmt.Println("connect redis success ...")
	defer rdb.Close()
}