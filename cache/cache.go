package cache

import (
	"context"
	"time"
)

// CacheClient is the interface that both Redis and Redis Cluster clients will implement.
type CacheClient interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string, value int64) (int64, error)
	SetAll(ctx context.Context, values map[string]interface{}) error
	GetAll(ctx context.Context, keys []string) (map[string]string, error)
	SetJson(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	SetJsonD(ctx context.Context, key string, value interface{}) error
}
