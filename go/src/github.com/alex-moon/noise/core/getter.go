package core

import (
    "github.com/garyburd/redigo/redis"
    "fmt"
)

type Getter interface {
    Get(key, member string) Item
}

type RedisGetter struct {
    conn redis.Conn
}

func NewGetter() RedisGetter {
    c, err := redis.Dial("tcp", Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("GETTER  -  Could not connect to Redis"))
    }

    return RedisGetter{c}
}

func (g RedisGetter) Get(key, member, default_val interface{}) Item {
    item, err := g.conn.Do("ZSCORE", key, member)
    if err != nil {
        panic(fmt.Sprintf("GETTER  -  could not get %s %s: %s\n", key, member, err.Error()))
    }
    if item == nil{
        return default_val
    }
    return item
}