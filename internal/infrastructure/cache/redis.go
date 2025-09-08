package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	Client *redis.Client
	ctx    context.Context
}

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisCache(config Config) (*RedisCache, error) {
	cache := &RedisCache{
		ctx: context.Background(),
	}

	cache.Client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	_, err := cache.Client.Ping(cache.ctx).Result()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// Set matches the interface signature
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if c.Client == nil {
		return nil
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Client.Set(c.ctx, key, jsonData, expiration).Err()
}

func (c *RedisCache) Get(key string, dest interface{}) error {
	if c.Client == nil {
		return redis.Nil
	}

	val, err := c.Client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Delete(key string) error {
	if c.Client == nil {
		return nil
	}
	return c.Client.Del(c.ctx, key).Err()
}

func (c *RedisCache) Exists(key string) bool {
	if c.Client == nil {
		return false
	}

	result, err := c.Client.Exists(c.ctx, key).Result()
	return err == nil && result > 0
}
