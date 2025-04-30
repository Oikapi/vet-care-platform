package redis

import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
)

type InventoryCache struct {
    client *redis.Client
}

func NewInventoryCache(addr string) *InventoryCache {
    client := redis.NewClient(&redis.Options{Addr: addr})
    return &InventoryCache{client: client}
}

func (c *InventoryCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, 0).Err()
}

func (c *InventoryCache) Get(ctx context.Context, key string) (string, error) {
    return c.client.Get(ctx, key).Result()
}