package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisClusterCacheClient is an implementation of the CacheClient interface for Redis Cluster
type RedisClusterCacheClient struct {
	client redis.UniversalClient
}

// NewRedisClusterCacheClient creates a new RedisClusterCacheClient instance for Redis Cluster
func NewRedisClusterCacheClient(addrs []string) *RedisClusterCacheClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs, // example: []string{"localhost:7000", "localhost:7001"}
	})
	return &RedisClusterCacheClient{client: rdb}
}

// Set stores a key-value pair in Redis
func (r *RedisClusterCacheClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

// Get retrieves a value from Redis by key
func (r *RedisClusterCacheClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del deletes a key from Redis
func (r *RedisClusterCacheClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisClusterCacheClient) Exists(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, key).Result()
	return res > 0, err
}

// Increment increments the integer value of a key by a given value
func (r *RedisClusterCacheClient) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// Decrement decrements the integer value of a key by a given value
func (r *RedisClusterCacheClient) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.DecrBy(ctx, key, value).Result()
}

// SetAll stores multiple key-value pairs in Redis
func (r *RedisClusterCacheClient) SetAll(ctx context.Context, values map[string]interface{}) error {
	pipe := r.client.Pipeline()
	for key, value := range values {
		pipe.Set(ctx, key, value, 0)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// GetAll retrieves multiple key-value pairs from Redis
func (r *RedisClusterCacheClient) GetAll(ctx context.Context, keys []string) ([]interface{}, error) {
	return r.client.MGet(ctx, keys...).Result()
}

// SetJson stores a JSON object in Redis with an expiration time
func (r *RedisClusterCacheClient) SetJson(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	// Marshal the value into JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// Store the JSON in Redis
	return r.client.Set(ctx, key, jsonData, expiry).Err()
}

// SetJsonD stores a JSON object in Redis with the default expiration time (0 means no expiry)
func (r *RedisClusterCacheClient) SetJsonD(ctx context.Context, key string, value interface{}) error {
	return r.SetJson(ctx, key, value, 0)
}
