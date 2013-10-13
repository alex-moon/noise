package main

import (
    "fmt"
    "time"
    // "github.com/fzzy/radix/redis"
    "github.com/garyburd/redigo/redis"
)

type MappedNGram struct {
    key string
    value string
}

func start_pinging(butt []byte) {
    one_second, err := time.ParseDuration("1s")
    if err != nil { /* whatever */ }

    for {
        fmt.Printf("%s...\n", butt)
        time.Sleep(one_second)
    }
}

func main() {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil { /* whatever */ }

    var pubsub = redis.PubSubConn{c}

    pubsub.Subscribe("noise")
    go func() { for {
        switch v := pubsub.Receive().(type) {
            case redis.Message:
                fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
                go start_pinging(v.Data)
        }
    } }()

    start_pinging([]byte("ping"))
}