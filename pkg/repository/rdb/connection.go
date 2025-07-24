package rdb

import (
	"context"
	"fmt"

	"github.com/kehl-gopher/logi/internal/config"
	"github.com/kehl-gopher/logi/internal/utils"
	"github.com/redis/go-redis/v9"
)

type RedisDB interface {
	ConnectRedis(ctx context.Context) error
	PingCache(ctx context.Context) bool
	Close()
}

type redisConn struct {
	log  *utils.Log
	conf *config.Config
	red  *redis.Client
}

func NewRedisConn(log *utils.Log, conf *config.Config) RedisDB {
	return &redisConn{log: log, conf: conf}
}

func (r *redisConn) ConnectRedis(ctx context.Context) error {
	port, err := utils.PortResolver(r.conf.RedisDB.Port)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", r.conf.RedisDB.Address, port)
	red := redis.NewClient(
		&redis.Options{
			Addr: addr,
			DB:   0,
		},
	)
	if err := red.Ping(ctx).Err(); err != nil {
		utils.PrintLog(r.log, err.Error(), utils.FatalLevel)
		return err
	}
	r.red = red
	utils.PrintLog(r.log, "redis connection successful", utils.InfoLevel)
	return nil
}

func (r *redisConn) PingCache(ctx context.Context) bool {
	return r.red.Ping(ctx).Err() == nil
}

func (r *redisConn) Close() {
	if r.red != nil {
		r.red.Close()
	}
}
