package redis

import (
	"blogfa/auth/config"
	zapLogger "blogfa/auth/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Storage   store = &rds{}
	logger    *zap.Logger
	namespace string = ""
	once      sync.Once
)

// store interface is interface for store things into redis
type store interface {
	Connect(config config.GlobalConfig) error
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key ...string) error
}

// rds struct for redis client
type rds struct {
	db *redis.Client
}

// Connect, method for connect to redis
func (r *rds) Connect(confs config.GlobalConfig) error {
	var err error

	once.Do(func() {
		logger = zapLogger.GetZapLogger(confs.Debug())

		r.db = redis.NewClient(&redis.Options{
			DB:       confs.Redis.DB,
			Addr:     confs.Redis.Host,
			Username: confs.Redis.Username,
			Password: confs.Redis.Password,
		})

		if err = r.db.Ping(context.Background()).Err(); err != nil {
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())
		}
	})

	return err
}

// Set meth a new key,value
func (r *rds) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}
	return r.db.Set(ctx, key, p, duration).Err()
}

// Get meth, get value with key
func (r *rds) Get(ctx context.Context, key string, dest interface{}) error {
	p, err := r.db.Get(ctx, key).Result()

	if p == "" {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return fmt.Errorf("Value Not Found")
	}

	if err != nil {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}

	return json.Unmarshal([]byte(p), &dest)
}

// Del for delete keys in redis
func (r *rds) Del(ctx context.Context, key ...string) error {
	_, err := r.db.Del(ctx, key...).Result()
	if err != nil {
		zapLogger.Prepare(logger).
			Development().
			Level(zap.ErrorLevel).
			Commit(err.Error())
		return err
	}
	return nil
}
