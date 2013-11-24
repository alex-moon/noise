package terms

import (
    "github.com/garyburd/redigo/redis"
    "github.com/alex-moon/noise/core"
    "fmt"
)

// TODO move to core
type Locker interface {
    Lock(key string)
    Unlock(key string)
    Create(key string)
}

type RedisTermLocker struct {
    conn redis.Conn
}

var term_locker *RedisTermLocker = nil

func TermLocker() RedisTermLocker {
    if term_locker == nil {
        c, err := redis.Dial("tcp", core.Config().Redis.Address)
        if err != nil {
            panic(fmt.Sprintf("TERM LOCKER  -  Could not connect to Redis"))
        }

        term_locker = &RedisTermLocker{c}
    }
    return *term_locker
}

func (l RedisTermLocker) Lock(term string) {
    l.conn.Do("BRPOP", core.Config().Mutex.Prefix + term)
}

func (l RedisTermLocker) Unlock(term string) {
    l.conn.Do("LPUSH", core.Config().Mutex.Prefix + term)
    recover() // not sure I'm using this right
}

func (l RedisTermLocker) Create(term string) {
    // TODO some way to tell if this has been created already... blah
    // if l.c.Do("RPOP", core.Config().Mutex.Prefix + term) != nil {
        l.conn.Do("LPUSH", core.Config().Mutex.Prefix + term)
    // }
}