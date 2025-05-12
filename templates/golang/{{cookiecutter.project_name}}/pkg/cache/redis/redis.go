// Package redis 提供 Redis 缓存服务
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/redis"
)

// Redis 缓存的 key 前缀
const cacheKeyPrefix = "{{cookiecutter.project_name}}"

// Cache redis 缓存
type Cache struct {
	name      string
	keyPrefix string
	cache     *cache.Cache
	ttl       time.Duration
}

// New 创建 redis 缓存实例
func New(name string, ttl time.Duration) *Cache {
	// key: {cache_key_prefix}:{cache_name}:{raw_key}
	keyPrefix := fmt.Sprintf("%s:%s", cacheKeyPrefix, name)

	c := cache.New(&cache.Options{
		Redis: redis.Client(),
		// Q：为什么不推荐利用 go-redis/cache 的本地缓存功能
		// A：缓存最好为单一 Backend，多 Backend 会给问题排查带来各种麻烦（除非业务真的有需求）
		LocalCache: nil,
	})
	return &Cache{
		name:      name,
		keyPrefix: keyPrefix,
		cache:     c,
		ttl:       ttl,
	}
}

// 生成实际 Redis 使用的 key
func (c *Cache) genKey(key string) string {
	return c.keyPrefix + ":" + key
}

// Set 将 value 存储到 redis 中，若 ttl 为 0 则使用默认值（Cache.ttl）
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if ttl == time.Duration(0) {
		ttl = c.ttl
	}

	item := cache.Item{Ctx: ctx, Key: c.genKey(key), Value: value, TTL: ttl}
	return c.cache.Set(&item)
}

// Exists 检查 key 在 redis 中是否存在
func (c *Cache) Exists(ctx context.Context, key string) bool {
	return c.cache.Exists(ctx, c.genKey(key))
}

// Get 从 redis 中获取值，并存储到 value 中，若失败则返回 error
func (c *Cache) Get(ctx context.Context, key string, value any) error {
	return c.cache.Get(ctx, c.genKey(key), value)
}

// Delete 从 redis 中删除指定的键
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, c.genKey(key))
}
