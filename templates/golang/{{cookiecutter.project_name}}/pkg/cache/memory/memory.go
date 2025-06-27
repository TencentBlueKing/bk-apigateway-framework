// Package memory 提供内存缓存服务（基于 freecache 封装，内存预分配 + LRU 算法）
// ref: https://github.com/coocood/freecache
package memory

import (
	"sync"

	"github.com/coocood/freecache"

	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
)

var (
	cache    *freecache.Cache
	initOnce sync.Once
)

// InitCache 根据指定容量初始化内存缓存（单位：MB）
func InitCache(capacity int) {
	initOnce.Do(func() {
		cache = freecache.NewCache(capacity * 1024 * 1024)
	})
}

// Cache 获取 cache 实例（提供 Get，Set，Del 等方法）
func Cache() *freecache.Cache {
	if cache == nil {
		log.Fatal("cache not init")
	}
	return cache
}
