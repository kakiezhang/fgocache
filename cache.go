package fgocache

import (
	"github.com/go-redis/redis"
	"github.com/kakiezhang/fgocache/hashring"
)

func New(servers []string) RedisCache {
	rc := RedisCache{
		nodes:   []string{},
		ring:    hashring.New(servers),
		clients: make(map[string]*redis.Client),
	}
	return rc
}
