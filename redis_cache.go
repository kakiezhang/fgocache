package fgocache

import (
	"errors"
	"fgocache/hashring"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

type RedisCache struct {
	nodes   []string
	ring    *(hashring.HashRing)
	clients map[string]*redis.Client
}

func (rc RedisCache) keyToConn(key string) *redis.Client {
	url, ok := rc.ring.GetNode(key)
	if ok {
		return rc.getConn(url)
	} else {
		panic(errors.New("node url format error"))
	}
}

func (rc RedisCache) getConn(url string) *redis.Client {
	cli, ok := rc.clients[url]

	if ok {
		return cli
	} else {
		arr := strings.Split(url, "/")
		if len(arr) != 4 {
			panic(errors.New("node url format error"))
		}

		addr := arr[2]
		db, err := strconv.Atoi(arr[3])
		if err != nil {
			panic(errors.New("node db must be integer"))
		}

		cli = NewRedisClient(addr, "", db)
		rc.clients[url] = cli
		return cli
	}
}

func (rc RedisCache) Get(key string) string {
	cli := rc.keyToConn(key)
	val, err := cli.Get(key).Result()

	defer func() {
		if p := recover(); p != nil {
			log.Printf("CMD[Get] panic: %v", p)
		}
	}()

	if err == redis.Nil {
		log.Printf("key[%v] not exist", key)
		return ""
	} else if err != nil {
		panic(err)
	} else {
		log.Printf("key[%v] => val[%v]", key, val)
		return val
	}
}

func (rc RedisCache) getExpiration(timeout int) time.Duration {
	if timeout <= 0 {
		timeout = 0
	}
	return time.Second * time.Duration(timeout)
}

func (rc RedisCache) Set(key, val string, timeout int) bool {
	exp := rc.getExpiration(timeout)

	cli := rc.keyToConn(key)
	err := cli.Set(key, val, exp).Err()

	defer func() {
		if p := recover(); p != nil {
			log.Printf("CMD[Set] panic: %v", p)
		}
	}()

	if err != nil {
		panic(err)
	}

	return true
}

func (rc RedisCache) Del(key string) bool {
	cli := rc.keyToConn(key)
	cnt, err := cli.Del(key).Result()

	defer func() {
		if p := recover(); p != nil {
			log.Printf("CMD[Del] panic: %v", p)
		}
	}()

	if err != nil {
		panic(err)
	}

	if cnt > 0 {
		log.Printf("key[%v] deleted", key)
		return true
	} else {
		return false
	}
}

func NewRedisClient(addr, password string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return cli
}
