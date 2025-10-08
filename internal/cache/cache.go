package cache

import (
	"net/http"
	"time"
)

const (
	HIT  = "HIT"
	MISS = "MISS"
)

type Cache struct {
	Response     *http.Response
	ResponseBody []byte
	CreatedAt    time.Time
}

func InitCache(r *http.Response, body []byte, created time.Time) *Cache {
	return &Cache{
		Response:     r,
		ResponseBody: body,
		CreatedAt:    created,
	}
}
