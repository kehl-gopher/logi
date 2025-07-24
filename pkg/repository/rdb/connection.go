package rdb

import (
	"fmt"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/redis/go-redis/v9"
)

type RedisDB interface {
	ConnectRedis() (*redis.Client, error)
}

type redisConn struct {
	log  *utils.Log
	conf *config.Config
}

func NewRedisConn(log *utils.Log, conf *config.Config) RedisDB {
	return &redisConn{log: log, conf: conf}
}

func (r *redisConn) ConnectRedis() (*redis.Client, error) {
	port, err := utils.PortResolver(r.conf.RedisDB.Port)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%d", r.conf.RedisDB.Address, port)
	red := redis.NewClient(
		&redis.Options{
			Addr: addr,
			DB:   0,
		},
	)
	return red, nil
}
