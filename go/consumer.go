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
        for {
            text_uuid, err := c.conn.Do("RPOP", c.channel)
            if err == nil && text_uuid != nil {
                uuid, ok := text_uuid.([]byte)
                if ok {
                    reader := NewReader(string(uuid))
                    term_counter := NewTermCounter(reader)
                    go term_counter.Run()
                } else {
                    panic(fmt.Sprintf("CONSUMER - Can't decode string %s\n", uuid))
                }
            } else {
                if err != nil {
                    panic(fmt.Sprintf("CONSUMER  -  Error: %s\n", err.Error()))
                }
                break
            }
        }
        <- c.subscriber.notifier
    }
}