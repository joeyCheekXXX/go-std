package captcha

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/mojocn/base64Captcha"
)

func NewDefaultRedisStore(_redis redis.UniversalClient) *RedisStore {
	return &RedisStore{
		redis:      _redis,
		Expiration: time.Second * 180,
		PreKey:     "CAPTCHA_",
		Context:    context.TODO(),
	}
}

type RedisStore struct {
	redis      redis.UniversalClient
	Expiration time.Duration
	PreKey     string
	Context    context.Context
}

func (rs *RedisStore) UseWithCtx(ctx context.Context) base64Captcha.Store {
	rs.Context = ctx
	return rs
}

func (rs *RedisStore) Set(id string, value string) error {
	err := rs.redis.Set(rs.Context, rs.PreKey+id, value, rs.Expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStore) Get(key string, clear bool) string {
	val, err := rs.redis.Get(rs.Context, key).Result()
	if err != nil {
		return ""
	}
	if clear {
		err := rs.redis.Del(rs.Context, key).Err()
		if err != nil {
			return ""
		}
	}
	return val
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	key := rs.PreKey + id
	v := rs.Get(key, clear)
	return v == answer
}
