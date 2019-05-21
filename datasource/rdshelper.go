package datasource

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"lottery/conf"
	"sync"
	"time"
)

var rdsLock sync.Mutex
var cacheInstance *RedisConn

// 封装成一个redis资源池
type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

// 对外只有一个命令，封装了一个redis的命令
func (rds *RedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := rds.pool.Get()
	defer conn.Close()

	t1 := time.Now().UnixNano()
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("rdshelper Do", err, e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug {
		fmt.Printf("[redis] [info] [%dus]cmd=%s, err=%s, args=%v, reply=%s\n", (t2-t1)/1000, commandName, err, args, reply)
	}
	return reply, err
}

/**
 * 分布式锁
script := `
local key   = KEYS[1]
local value = ARGV[1]
local ttl   = ARGV[2]

local ok = redis.call('setnx', key, value)
if ok == 1 then
  redis.call('expire', key, ttl)
end

return ok
`
 */
func (rds *RedisConn) RedisLock(script string, args ...interface{}) (reply interface{}, err error) {
	c := rds.pool.Get()
	var getScript = redis.NewScript(1, script)
	reply, err = getScript.Do(c, args...)

	if err != nil {
		return nil, err
	}
	return reply, err
}

// 设置是否打印操作日志
func (rds *RedisConn) ShowDebug(b bool) {
	rds.showDebug = b
}

// 得到唯一的redis缓存实例
func InstanceCache() *RedisConn {
	if cacheInstance != nil {
		return cacheInstance
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()

	if cacheInstance != nil {
		return cacheInstance
	}
	return NewCache()
}

// 重新实例化
func NewCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.RdsCache.Host, conf.RdsCache.Port))
			if err != nil {
				log.Fatal("rdshelper.NewCache Dial error ", err)
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000,
		MaxActive:       10000,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	instance := &RedisConn{
		pool: &pool,
	}
	cacheInstance = instance
	cacheInstance.ShowDebug(false)
	//cacheInstance.ShowDebug(false)
	return cacheInstance
}
