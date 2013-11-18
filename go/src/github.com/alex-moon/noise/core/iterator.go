package core

import (
    "io/ioutil"
    "github.com/garyburd/redigo/redis"
)

type Iterator interface {
    Items() chan string
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

func (i QueueIterator) Items() chan string {
    items := make(chan string)
    go func() {
        for item, err := i.conn.Do("RPOP", i.channel) {
            if err == nil && item != nil {
                items <- string(item)
            } else {
                if err != nil {
                    panic(fmt.Sprintf("QUEUE ITERATOR %s  -  RPOP error from Redis: %s", i.channel, err.Error()))
                }
                items <- nil
                break
            }
        }
    }
    return items
}


// REDIS SORTED SET

type SetIterator struct {
    conn redis.Conn
    key string
    min float
    max float
}

func NewSetIterator(key string) SetIterator {
    return NewSetIterator(key, 0, -1)
}

func NewSetIterator(key string, min float, max float) SetIterator {
    c, err := redis.Dial("tcp", Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("SET ITERATOR %s  -  Could not connect to Redis", key))
    }

    return SetIterator {
        conn: c,
        key: key,
        min: min,
        max: max,
    }
}

// TODO: this violates interface Iterator above - nice way to do polymorphism in Go?
// alternatively: a Reader for these which just gets the score - however, that does create a race condition which I don't like at all
func (i SetIterator) Items() chan string, float {
    items := make(chan string, float)
    go func() {
        
    }
}

// FILESYSTEM

type FileSystemIterator struct {
    dir string
}

func NewFileSystemIterator(dir) {
    return FileSystemIterator{dir}
}

func (i FileSystemIterator) Items() chan string {
    files, err := ioutil.ReadDir(i.dir)
    if err != nil {
        panic(fmt.Sprintf("FILESYSTEM ITERATOR %s - could not read dir: %s", i.dir, err.Error()))
    }
    items := make(chan string)
    go func() {
        for _, filestat := range files {
            filename := filestat.Name()
            items <- i.dir + "/" + string(filename)
        }
        items <- nil
    }
    return items
}