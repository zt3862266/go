package redis

import (
	"errors"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/zt3862266/go/config"
	. "github.com/zt3862266/go/log"
)

var (
	IdleTimeout = time.Minute
	cachePool   []*redis.Pool
)

type RongCache struct {
}

func InitRedis() {

	redisConf := config.GlobalEnv.Redis
	for _, conf := range redisConf.Machine {
		Pool := &redis.Pool{
			MaxIdle:     redisConf.MaxIdle,
			MaxActive:   redisConf.MaxActive,
			IdleTimeout: IdleTimeout,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", conf.Addr,
					redis.DialReadTimeout(time.Second*time.Duration(redisConf.ReadTimeoutSeconds)),
					redis.DialWriteTimeout(time.Second*time.Duration(redisConf.WriteTimeoutSeconds)),
					redis.DialConnectTimeout(time.Second*1),
					redis.DialKeepAlive(time.Minute*5),
				)
				if err != nil {
					return nil, err
				}
				if conf.Password != "" {
					if _, err := c.Do("AUTH", conf.Password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
		}
		cachePool = append(cachePool, Pool)
		Info("cachePool:%v", cachePool)
	}
}

func (r *RongCache) getPoll() (*redis.Pool, error) {
	poolLen := len(cachePool)
	if poolLen == 0 {
		return nil, errors.New("no valid slave")
	}
	chooseIdx := rand.Intn(poolLen)
	return cachePool[chooseIdx], nil

}

func (r *RongCache) Set(key, value string, ttl int) (err error) {
	pool, err := r.getPoll()
	if err != nil {
		return err
	}
	conn := pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value, "EX", ttl)
	if err != nil {
		Error("set failed:%v", err.Error())

	}
	return err

}

func (r *RongCache) Get(key string) (value string, err error) {
	pool, err := r.getPoll()
	if err != nil {
		return "", err
	}
	conn := pool.Get()
	defer conn.Close()
	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		Error("get failed:%v", err.Error())
	}
	return value, err
}
