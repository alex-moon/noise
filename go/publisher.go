package main

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

type Publisher struct {
    channel string
    conn redis.Conn
}

func NewPublisher(channel string) Publisher {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("PUBLISHER %s  -  Could not connect to Redis", channel))
    }

    return Publisher{
        conn: c,
        channel: channel,
    }
}

func (p Publisher) Publish(channel string, value string) {
    p.conn.Do("LPUSH", channel, value)
    p.conn.Do("PUBLISH", channel, 1)
}