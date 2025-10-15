package cache

import (
	"net/http"
	"sync"
	"time"
)

const (
	CACHE_HIT  = "HIT"
	CACHE_MISS = "MISS"
)

type Cache struct {
	CacheObjs map[string]*CacheObject
	mutex     sync.RWMutex
}

func InitCache() *Cache {
	return &Cache{
		CacheObjs: make(map[string]*CacheObject),
		mutex:     sync.RWMutex{},
	}
}

type CacheObject struct {
	Response     *http.Response
	ResponseBody []byte
	CreatedAt    time.Time
}

func CreateCacheObject(r *http.Response, body []byte, createdAt time.Time) *CacheObject {
	return &CacheObject{
		Response:     r,
		ResponseBody: body,
		CreatedAt:    createdAt,
	}
}

func (c *Cache) Get(key string) (*CacheObject, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if val, ok := c.CacheObjs[key]; ok {
		return val, true
	}
	return nil, false
}

func (c *Cache) Set(key string, cache *CacheObject) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.CacheObjs[key] = cache
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.CacheObjs, key)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k := range c.CacheObjs {
		delete(c.CacheObjs, k)
	}
}
