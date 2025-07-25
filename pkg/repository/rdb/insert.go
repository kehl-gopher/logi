package rdb

import (
	"context"
	"fmt"
	"time"
)

func (rdb *redisConn) Set(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	ins := rdb.red.Set(ctx, key, data, ttl)
	if ins.Err() != nil {
		return fmt.Errorf("failed to set redis key %s: %w", key, ins.Err())
	}
	return nil
}

func (rdb *redisConn) Queue(ctx context.Context, key string, data interface{}) error {
	err := rdb.red.LPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("unable to add %s to queue: %w", key, err)
	}
	return nil
}
