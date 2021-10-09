package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var localCache *cache.Cache

func init() {
	localCache = cache.New(5*time.Minute, 10*time.Minute)
	if localCache == nil {
		panic("创建本地缓存对象为空")
	}
}

func get(key string) (interface{}, bool) {
	return localCache.Get(key)
}

func set(key string, object interface{}, duration time.Duration) {
	localCache.Set(key, object, duration)
}

func delete(key string) {
	localCache.Delete(key)
	localCache.DeleteExpired()
}
