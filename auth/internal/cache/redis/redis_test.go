package redis

import (
	"auth/internal/config"
	"auth/internal/model/principal"
	"auth/pkg/logger"
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.RedisConfig
		l   logger.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisCache
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCache_Get(t *testing.T) {
	type fields struct {
		rdb      *redis.Client
		expireIn time.Duration
		l        logger.Logger
	}

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantP   principal.Principal
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				rdb:      tt.fields.rdb,
				expireIn: tt.fields.expireIn,
				l:        tt.fields.l,
			}
			gotP, err := r.Get(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("RedisCache.Get() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

func TestRedisCache_Set(t *testing.T) {
	type fields struct {
		rdb      *redis.Client
		expireIn time.Duration
		l        logger.Logger
	}
	type args struct {
		email string
		p     principal.Principal
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				rdb:      tt.fields.rdb,
				expireIn: tt.fields.expireIn,
				l:        tt.fields.l,
			}
			if err := r.Set(tt.args.email, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("RedisCache.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisCache_Del(t *testing.T) {
	type fields struct {
		rdb      *redis.Client
		expireIn time.Duration
		l        logger.Logger
	}
	type args struct {
		email string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAffected bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				rdb:      tt.fields.rdb,
				expireIn: tt.fields.expireIn,
				l:        tt.fields.l,
			}
			if gotAffected := r.Del(tt.args.email); gotAffected != tt.wantAffected {
				t.Errorf("RedisCache.Del() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}

func TestRedisCache_EntryCount(t *testing.T) {
	type fields struct {
		rdb      *redis.Client
		expireIn time.Duration
		l        logger.Logger
	}
	tests := []struct {
		name           string
		fields         fields
		wantEntryCount int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				rdb:      tt.fields.rdb,
				expireIn: tt.fields.expireIn,
				l:        tt.fields.l,
			}
			if gotEntryCount := r.EntryCount(); gotEntryCount != tt.wantEntryCount {
				t.Errorf("RedisCache.EntryCount() = %v, want %v", gotEntryCount, tt.wantEntryCount)
			}
		})
	}
}

func TestRedisCache_HitCount(t *testing.T) {
	type fields struct {
		rdb      *redis.Client
		expireIn time.Duration
		l        logger.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				rdb:      tt.fields.rdb,
				expireIn: tt.fields.expireIn,
				l:        tt.fields.l,
			}
			if got := r.HitCount(); got != tt.want {
				t.Errorf("RedisCache.HitCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCache_MissCount(t *testing.T) {
	type fields struct {
		rdb      *redis.Client
		expireIn time.Duration
		l        logger.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				rdb:      tt.fields.rdb,
				expireIn: tt.fields.expireIn,
				l:        tt.fields.l,
			}
			if got := r.MissCount(); got != tt.want {
				t.Errorf("RedisCache.MissCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testcontainerRedis(cfg *config.RedisConfig) (*RedisCache, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis/redis-stack-server:6.2.2-v5",
		ExposedPorts: []string{"6379/tcp"},
		Env: map[string]string{
			"REDIS_ARGS": "--requirepass redis-password",
		},
		WaitingFor: wait.ForListeningPort("6379/tcp"),
	}

	container, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(context.TODO())
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(context.TODO(), "6379")
	if err != nil {
		return nil, err
	}

	cfg.Host = ip
	cfg.Port = mappedPort.Port()

	return New(cfg, logger.New(logger.Debug))
}
