package core

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

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

func (g RedisScoreGetter) GetFloat(key, member string, default_val float64) float64 {
    item, err := g.conn.Do("ZSCORE", key, member)
    if err != nil {
        fmt.Sprintf("GETTER  -  could not get %s %s - got %s: %s - retrying...\n", key, member, item, err.Error())
        return g.GetFloat(key, member, default_val)
    }
    score, err := redis.Float64(item, err)
    if err == redis.ErrNil {
        return default_val
    } else if err != nil {
        fmt.Sprintf("GETTER  -  could not convert %s to float64: %s - retrying...\n", item, err.Error())
        return g.GetFloat(key, member, default_val)
    }
    return score
}


func (g RedisScoreGetter) GetInt(key, member string, default_val int) int {
    score_float := g.GetFloat(key, member, float64(default_val))
    return int(score_float)
}