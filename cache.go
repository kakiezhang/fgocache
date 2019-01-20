package fgocache

import (
	"fgocache/hashring"
	"github.com/go-redis/redis"
)

func New(servers []string) RedisCache {
	rc := RedisCache{
		nodes:   []string{},
		ring:    hashring.New(servers),
		clients: make(map[string]*redis.Client),
	}
	return rc
}
