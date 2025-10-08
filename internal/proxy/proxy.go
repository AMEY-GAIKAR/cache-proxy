package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"cache-proxy/internal/cache"
)

type Proxy struct {
	Origin string
	Cache  map[string]*cache.Cache
	Mutex  sync.RWMutex
}

func InitProxy(origin string) *Proxy {
	return &Proxy{
		Origin: origin,
		Cache:  make(map[string]*cache.Cache),
		Mutex:  sync.RWMutex{},
	}
}

func (p *Proxy) ClearCache() {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.Cache = make(map[string]*cache.Cache)
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

	p.Mutex.RLock()
	if val, ok := p.Cache[key]; ok {
		p.Mutex.RUnlock()
		WriteResponseWithHeaders(w, val.Response, val.ResponseBody, cache.HIT, key)
		return
	}
	p.Mutex.RUnlock()

	log.Printf("Cache not present for key: %s\n", key)
	resp, err := http.Get(p.Origin + r.URL.String())
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading content", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	p.Mutex.Lock()
	p.Cache[key] = cache.InitCache(resp, body, time.Now())
	p.Mutex.Unlock()

	WriteResponseWithHeaders(w, resp, body, cache.MISS, key)
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
