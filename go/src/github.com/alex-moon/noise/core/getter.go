package core

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

type Getter interface {
    Get(key, member string) Item
}

type RedisScoreGetter struct {
    conn redis.Conn
}

func NewGetter() RedisScoreGetter {
    c, err := redis.Dial("tcp", Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("GETTER  -  Could not connect to Redis"))
    }

    return RedisScoreGetter{c}
}

func (g RedisScoreGetter) Get(key, member, default_val interface{}) Item {
    item, err := redis.Float64(g.conn.Do("ZSCORE", key, member))
    if err == redis.ErrNil {
        return default_val
    } else if err != nil {
        panic(fmt.Sprintf("GETTER  -  could not get %s %s: %s\n", key, member, err.Error()))
    }
    return float32(item)
}