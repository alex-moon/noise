package main

import (
    "fmt"
    "strconv"
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
        msg, err := s.conn.Receive()
        if err == nil {
            msg, ok := msg.(string)
            if !ok { panic(msg) }
            number_of_texts, err := strconv.Atoi(msg)
            if err != nil { panic(err) }
            s.notifier <- number_of_texts
        } else {
            panic(fmt.Sprintf("SUBSCRIBER %s %s  -  Receive error: %s", s.channel, s.notifier, err.Error()))
        }
    }
}