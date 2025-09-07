package cache

import (
	"context"
	"encoding/json"
	"time"

	"dashboard-starter/pkg/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Init(config Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		logger.Warn("Redis connection failed, running without cache", zap.Error(err))
		return err
	}

	logger.Info("Redis cache connected successfully")
	return nil
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if Client == nil {
		return nil // Graceful degradation
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return Client.Set(ctx, key, jsonData, expiration).Err()
}

func Get(key string, dest interface{}) error {
	if Client == nil {
		return redis.Nil // Graceful degradation
	}

	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func Delete(key string) error {
	if Client == nil {
		return nil
	}
	return Client.Del(ctx, key).Err()
}

func Exists(key string) bool {
	if Client == nil {
		return false
	}

	result, err := Client.Exists(ctx, key).Result()
	return err == nil && result > 0
}
