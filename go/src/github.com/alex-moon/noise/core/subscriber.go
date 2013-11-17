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
    pubsub := redis.PubSubConn{Conn: s.conn}
    pubsub.Subscribe(s.channel)

    for {
        switch data := pubsub.Receive().(type) {
        case redis.Message:
            msg := string(data.Data)
            number_of_texts, err := strconv.Atoi(msg)
            if err == nil {
                s.notifier <- number_of_texts
            } else { panic(fmt.Sprintf("SUBSCRIBER - could not convert %s to int: %s\n", msg, err.Error())) }
        case error:
            panic(fmt.Sprintf("SUBSCRIBER %s %s  -  Receive error: %s", s.channel, s.notifier, data))
        }
    }
}