// Package redis 提供了 Redis 相关的封装（基于 redis/go-redis/v9）
// SaaS 开发者查阅该文档以了解使用方法：https://redis.uptrace.dev/guide/go-redis.html
package redis

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/config"
	log "bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/logging"
)

var (
	rds      *redis.Client
	initOnce sync.Once
)

const (
	// 尝试连接超时 单位：s
	dialTimeout = 2
	// 读超时 单位：s
	readTimeout = 1
	// 写超时 单位：s
	writeTimeout = 1
	// 闲置超时 单位: s
	idleTimeout = 3 * 60
	// 连接池大小 / 核
	poolSizeMultiple = 20
	// 最小空闲连接数 / 核
	minIdleConnectionMultiple = 10
)

// InitRedisClient init redis client with config.RedisConfig
func InitRedisClient(ctx context.Context, cfg *config.RedisConfig) {
	opts, err := redis.ParseURL(cfg.DSN())
	if err != nil {
		log.Fatalf("redis parse url error: %s", err.Error())
	}

	// Redis 配置
	opts.DialTimeout = time.Duration(dialTimeout) * time.Second
	opts.ReadTimeout = time.Duration(readTimeout) * time.Second
	opts.WriteTimeout = time.Duration(writeTimeout) * time.Second
	opts.ConnMaxIdleTime = time.Duration(idleTimeout) * time.Second
	opts.PoolSize = poolSizeMultiple * runtime.NumCPU()
	opts.MinIdleConns = minIdleConnectionMultiple * runtime.NumCPU()

	initOnce.Do(func() {
		rds = redis.NewClient(opts)
		if _, err = rds.Ping(ctx).Result(); err != nil {
			log.Fatalf("redis connect error: %s", err.Error())
		} else {
			log.Infof(ctx, "redis: %s:%d/%d connected", cfg.Host, cfg.Port, cfg.DB)
		}
		// OpenTelemetry Tracing
		if err = redisotel.InstrumentTracing(rds); err != nil {
			log.Fatalf("failed to enable redis tracing instrumentation: %s", err)
		}
	})
}

// Client 获取 redis 客户端
func Client() *redis.Client {
	if rds == nil {
		log.Fatal("redis client not init")
	}
	return rds
}
