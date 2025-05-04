package redis

import (
    "context"
    "time"
    "encoding/json"
    "github.com/go-redis/redis/v8"
)

type InventoryCache struct {
    client *redis.Client
}

func NewInventoryCache(addr string) (*InventoryCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    // Проверяем подключение
    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }
    return &InventoryCache{client: client}, nil
}

func (c *InventoryCache) Set(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    // Устанавливаем TTL (например, 1 час)
    return c.client.Set(ctx, key, data, 3600*time.Second).Err()
}

func (c *InventoryCache) Get(ctx context.Context, key string) (interface{}, error) {
    val, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Ключ не найден, возвращаем nil без ошибки
    }
    if err != nil {
        return nil, err
    }
    return val, nil
}

func (c *InventoryCache) Del(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}