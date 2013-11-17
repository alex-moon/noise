package core

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
    c, err := redis.Dial("tcp", Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("CONSUMER %s  -  Could not connect to Redis", channel))
    }

    notifier := make(chan int)
    subscriber := NewSubscriber(channel, notifier)

    return Consumer {
        conn: c,
        channel: channel,
        subscriber: subscriber,
    }
}

func (c Consumer) Consume(new_iterator func() chan interface{}, work func(interface{})) {
    go c.subscriber.Subscribe()

    for {
        iterator := new_iterator()
        for item := range iterator {
            if item == nil { break }
            go work(item)
        }
        /* WAS:
        reader := NewReader()
        for text := range reader.texts {
            if text == nil { break }
            term_counter := NewTermCounter(text)
            go term_counter.Run()
        }
        */
        <- c.subscriber.notifier
    }
}