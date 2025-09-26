package ctxcache

import (
	"context"
	"sync"
)

type ctxCacheKey struct{}

// Init 初始化context的缓存, 使用sync.Map, 并发安全
func Init(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxCacheKey{}, new(sync.Map))
}

// Get 获取通过key获取指定类型的值
func Get[T any](ctx context.Context, key any) (value T, ok bool) {
	var zero T

	cacheMap, valid := ctx.Value(ctxCacheKey{}).(*sync.Map)
	if !valid {
		return zero, false
	}

	loadedValue, exists := cacheMap.Load(key)
	if !exists {
		return zero, false
	}

	if v, match := loadedValue.(T); match {
		return v, true
	}

	return zero, false
}

// Store 存储一个值
func Store(ctx context.Context, key any, obj any) {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		cacheMap.Store(key, obj)
	}
}

// HasKey 判断是否存在这个Key
func HasKey(ctx context.Context, key any) bool {
	if cacheMap, ok := ctx.Value(ctxCacheKey{}).(*sync.Map); ok {
		_, ok := cacheMap.Load(key)
		return ok
	}

	return false
}
