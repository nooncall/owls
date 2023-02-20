package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/db_info"
	"github.com/nooncall/owls/go/utils/logger"
)

func NewRedisCli(clusterName string, database int) (r *redis.Client, err error) {
	cluster, err := db_info.GetClusterByName(clusterName)
	if err != nil {
		return nil, fmt.Errorf("get cluster info err: %s", err.Error())
	}

	// new，and init
	redisCli := redis.NewClient(&redis.Options{
		Addr:        cluster.Addr,
		Password:    cluster.Pwd, // redis的认证密码
		DB:          database,    // 连接的database库
		IdleTimeout: 300,         // 默认Idle超时时间
		PoolSize:    10,          // 连接池
	})

	res, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		logger.Errorf("Connect Failed! Err: %v", err)
		return nil, err
	}
	logger.Errorf("Connect Successful! Ping => %v", res)
	return redisCli, nil
}
