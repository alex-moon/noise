package core

import (
    "fmt"
    "strconv"
    "io/ioutil"
    "github.com/garyburd/redigo/redis"
)

const SET_SCORE_ITERATOR int = 1
const SET_RANK_ITERATOR int = 2

type Item interface {}

type Iterator interface {
    Items() chan Item
}


// REDIS LIST

type QueueIterator struct {
    conn redis.Conn
    channel string
}

func NewQueueIterator(channel string) QueueIterator {
    c, err := redis.Dial("tcp", Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("QUEUE ITERATOR %s  -  Could not connect to Redis", channel))
    }
    
    return QueueIterator {
        conn: c,
        channel: channel,
    }
}

func (i QueueIterator) Items() chan Item {
    items := make(chan Item)
    go func() {
        for {
            item, err := i.conn.Do("RPOP", i.channel)
            if err == nil && item != nil {
                items <- item
            } else {
                if err != nil {
                    panic(fmt.Sprintf("QUEUE ITERATOR %s  -  RPOP error from Redis: %s", i.channel, err.Error()))
                }
                items <- nil
                break
            }
        }
    }()
    return items
}


// REDIS SORTED SET

type SetIterator struct {
    conn redis.Conn
    key string
    iterator_type int  // TODO really? Flags on native objects? Doesn't feel right to me...
}

func NewSetIterator(key string, iterator_type int) SetIterator {
    c, err := redis.Dial("tcp", Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("SET ITERATOR %s  -  Could not connect to Redis", key))
    }

    if iterator_type != SET_RANK_ITERATOR && iterator_type != SET_SCORE_ITERATOR {
        panic(fmt.Sprintf("SET ITERATOR %s  -  Did not recognise iterator type %s\n", key, iterator_type))
    }

    return SetIterator {
        conn: c,
        key: key,
        iterator_type: iterator_type,
    }
}

type SetMember struct {
    Term string
    Score float64
    SumTotal float64
}

func (set SetIterator) Items() chan Item {
    getter := NewGetter()
    items := make(chan Item)
    go func() {
        var min float64 = 0.0
        var max float64 = 1.0
        var sum_total float64 = 1.0

        switch set.iterator_type {
            case SET_RANK_ITERATOR:
                max = -1.0
                sum_total := getter.GetFloat(set.key, SET_SUM_MEMBER, 0.0)
                if sum_total == 0.0 {
                    panic(fmt.Sprintf("SET ITERATOR %s  -  Could not get sum total for set %s\n", set.key, set.iterator_type))
                }
            case SET_SCORE_ITERATOR:
                break
            default:
                panic(fmt.Sprintf("SET ITERATOR %s  -  Did not recognise iterator type %s\n", set.key, set.iterator_type))
        }

        members, err := redis.Strings(set.conn.Do("ZRANGE", set.key, min, max, "WITHSCORES"))
        if err != nil {
            panic(fmt.Sprintf("SET ITERATOR %s - Could not do ZRANGE %s", set.key, err.Error()))
        }

        for i := 0; i < len(members); i += 2 {
            member := members[i]
            if member == SET_SUM_MEMBER { continue }
            score, err := strconv.ParseFloat(members[i+1], 64)
            if err != nil {
                panic(fmt.Sprintf("SET ITERATOR %s - Could not convert %s to float %s", set.key, members[i+1], err.Error()))
            }
            items <- SetMember{member, score, sum_total}
        }

        items <- nil
    }()
    return items
}


// FILESYSTEM

type FileSystemIterator struct {
    dir string
}

func NewFileSystemIterator(dir string) FileSystemIterator {
    return FileSystemIterator{dir}
}

func (i FileSystemIterator) Items() chan Item {
    files, err := ioutil.ReadDir(i.dir)
    if err != nil {
        panic(fmt.Sprintf("FILESYSTEM ITERATOR %s - could not read dir: %s", i.dir, err.Error()))
    }
    items := make(chan Item)
    go func() {
        for _, filestat := range files {
            items <- string(filestat.Name())
        }
        items <- nil
    }()
    return items
}