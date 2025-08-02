package rdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (r *redisConn) Get(ctx context.Context, key string, receiver interface{}) (interface{}, error) {
	val, err := r.RED().Get(ctx, key).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("cannot fetch data from redis")
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), receiver)

	if err != nil {
		return nil, err
	}
	return receiver, nil
}
