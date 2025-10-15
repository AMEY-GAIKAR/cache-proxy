Cache Proxy Server 

In-Memory Cache Proxy Server in Go

Usage

```bash
go run cmd/main.go -port <int> -origin <url string> -clear--cache false
```

Example usage
```bash
go run cmd/main.go -port 8080 -origin http://localhost:3000 

```

TODO
- Cache TTL
- LRU & LFU eviction policies
- Redis support
