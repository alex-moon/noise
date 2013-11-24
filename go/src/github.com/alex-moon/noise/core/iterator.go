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
    Score float32
    SumTotal float32
    CrossReference []SetMember
}

func (set SetIterator) Items() chan Item {
    items := make(chan Item)
    go func() {
        var min float32 = 0.0
        var max float32 = 1.0
        var sum_total float32 = 1.0

        switch set.iterator_type {
            case SET_RANK_ITERATOR:
                max = -1.0
                t, err := redis.Float64(set.conn.Do("ZSCORE", set.key, SET_SUM_MEMBER))
                if err != nil {
                    panic(fmt.Sprintf("SET ITERATOR %s  -  Could not get sum total for set %s: %s\n", set.key, set.iterator_type, err.Error()))
                }
                sum_total = float32(t)
            case SET_SCORE_ITERATOR:
                break
            default:
                panic(fmt.Sprintf("SET ITERATOR %s  -  Did not recognise iterator type %s\n", set.key, set.iterator_type))
        }

        members, err := redis.Strings(set.conn.Do("ZRANGE", set.key, min, max, "WITHSCORES"))
        if err != nil {
            panic(fmt.Sprintf("SET ITERATOR %s - Could not do ZRANGE %s", set.key, err.Error()))
        }

        // first construct the cross-reference list we'll attach to each member (popping from the front each time)
        cross_reference := []SetMember{}
        for i := 2; i < len(members); i += 2 {
            cross_reference_term := members[i]
            cross_reference_score, err := strconv.ParseFloat(members[i+1], 32)
            if err != nil {
                panic(fmt.Sprintf("SET ITERATOR %s - Could not convert %s to float %s", set.key, members[i+1], err.Error()))
            }
            cross_reference_member := SetMember{cross_reference_term, float32(cross_reference_score), sum_total, []SetMember{}}
            cross_reference = append(cross_reference, cross_reference_member)
        }

        for i := 0; i < len(members); i += 2 {
            member := members[i]
            score, err := strconv.ParseFloat(members[i+1], 32)
            if err != nil {
                panic(fmt.Sprintf("SET ITERATOR %s - Could not convert %s to float %s", set.key, members[i+1], err.Error()))
            }
            cross_reference := cross_reference[1:]
            items <- SetMember{member, float32(score), sum_total, cross_reference}
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