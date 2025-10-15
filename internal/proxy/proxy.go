package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/AMEY-GAIKAR/cache-proxy/internal/cache"
)

type Proxy struct {
	Origin string
	Cache  *cache.Cache
	Mutex  sync.RWMutex
}

func InitProxy(originURL string) *Proxy {
	return &Proxy{
		Origin: originURL,
		Cache:  cache.InitCache(),
		Mutex:  sync.RWMutex{},
	}
}

func (p *Proxy) ClearCache() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.Cache.Clear()
	log.Printf("Cleared cache")
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/clear-cache" {
		log.Printf("Clearing cache...\n")
		p.ClearCache()
		w.Write([]byte("Cache cleared successfully"))
		return
	}

	key := fmt.Sprintf("%s:%s", r.Method, r.URL.Path)

	if val, ok := p.Cache.Get(key); ok {
		WriteResponseWithHeaders(w, val.Response, val.ResponseBody, cache.CACHE_HIT, key)
		return
	}

	log.Printf("Cache not present for: %s\n", key)

	resp, err := http.Get(p.Origin + r.URL.String())
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading content", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	newObj := cache.CreateCacheObject(resp, body, time.Now())
	p.Cache.Set(key, newObj)

	WriteResponseWithHeaders(w, resp, body, cache.CACHE_MISS, key)
}

func WriteResponseWithHeaders(w http.ResponseWriter, r *http.Response, body []byte, cacheHeader string, key string) {
	log.Printf("Cache: %s %s\n", cacheHeader, key)
	w.Header().Set("X-Cache", cacheHeader)
	w.WriteHeader(r.StatusCode)
	for k, v := range r.Header {
		w.Header()[k] = v
	}
	w.Write(body)
}
