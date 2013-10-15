package main

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

type Subscriber struct {
    channel string
    notifier chan int
    conn redis.Conn
}

func NewSubscriber(channel string, notifier chan int) Subscriber {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("SUBSCRIBER %s %s  -  Could not connect to Redis", channel, notifier))
    }

    return Subscriber{
        conn: c,
        channel: channel,
        notifier: notifier,
    }
}

func (s Subscriber) Subscribe() {
    s.conn.Do("SUBSCRIBE", s.channel)

    for {
        _, err := s.conn.Receive()
        if err == nil {
            s.notifier <- 1
        } else {
            panic(fmt.Sprintf("SUBSCRIBER %s %s  -  Receive error: %s", s.channel, s.notifier, err.Error()))
        }
    }
}