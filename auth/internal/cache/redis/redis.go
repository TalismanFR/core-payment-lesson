package redis

import (
	"auth/internal/config"
	"auth/internal/model/principal"
	"auth/pkg/helpers"
	"auth/pkg/logger"
	"context"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

type RedisCache struct {
	rdb      *redis.Client
	expireIn time.Duration
	l        logger.Logger
}

func New(cfg *config.RedisConfig, l logger.Logger) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{
		rdb:      rdb,
		expireIn: cfg.ExpireIn,
		l:        l,
	}, nil
}

func (r *RedisCache) Get(email string) (p principal.Principal, err error) {
	err = r.rdb.Get(context.TODO(), email).Scan(&p)
	return
}

func (r *RedisCache) Set(email string, p principal.Principal) error {
	return r.rdb.SetEX(context.TODO(), email, &p, r.expireIn).Err()
}

func (r *RedisCache) Del(email string) (affected bool) {
	res, err := r.rdb.Del(context.TODO(), email).Result()
	if err != nil {
		return false
	}
	if res > 0 {
		return true
	}
	r.l.Debug("RedisCache.Del: items deleted: %v", res)
	return false
}

func (r *RedisCache) EntryCount() (entryCount int64) {
	return r.rdb.DBSize(context.TODO()).Val()
}

func (r *RedisCache) HitCount() int64 {
	res, err := r.rdb.Info(context.TODO(), "stats").Bytes()
	if err != nil {
		return -1
	}

	hits, err := helpers.ParseField(res, []byte("keyspace_hits:"), "\n\r\n\r")
	if err != nil {
		r.l.Error("RedisCache.parseField: err: %v", err)
	}

	return hits
}

func (r *RedisCache) MissCount() int64 {
	res, err := r.rdb.Info(context.TODO(), "stats").Bytes()
	if err != nil {
		return -1
	}

	hits, err := helpers.ParseField(res, []byte("keyspace_misses:"), "\n\r\n\r")
	if err != nil {
		r.l.Error("RedisCache.parseField: err: %v", err)
	}

	return hits
}
