package main

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

type Consumer struct {
    channel string
    subscriber Subscriber
    conn redis.Conn
}

func NewConsumer(channel string) Consumer {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("CONSUMER %s  -  Could not connect to Redis", channel))
    }

    notifier := make(chan int)
    subscriber := NewSubscriber(channel, notifier)

    return Consumer{
        conn: c,
        channel: channel,
        subscriber: subscriber,
    }
}

func (c Consumer) Consume() {
    go c.subscriber.Subscribe()

    for {
        reader := NewReader()
        for text := range reader.texts {
            if text == nil { break }
            term_counter := NewTermCounter(text)
            go term_counter.Run()
        }
        <- c.subscriber.notifier
    }
}