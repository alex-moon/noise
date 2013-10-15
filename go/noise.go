package main

import (
    "fmt"
    "time"
    "github.com/garyburd/redigo/redis"
)

type MappedNGram struct {
    key string
    value string
}

func subscribe(channel string, notify chan int) {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil { panic("SUBSCRIBER  -  Could not connect to Redis") }

    fmt.Printf("SUBSCRIBER  -  About to subscribe...\n")
    c.Do("SUBSCRIBE", channel)
    fmt.Printf("SUBSCRIBER  -  subscribed!\n")

    for {
        _, err := c.Receive()
        if err == nil {
            fmt.Printf("\nSUBSCRIBER  - notify: new list elems to pop\n")
            notify <- 1
        } else {
            panic("SUBSCRIBER  -  Receive error: " + err.Error() + "\n")
        }
    }
}

func consume(channel string) {
    notify := make(chan int)
    go subscribe(channel, notify)

    c, err := redis.Dial("tcp", ":6379")
    if err != nil { panic("CONSUMER  -  Could not connect to Redis") }
    for {
        for {
            fmt.Print("CONSUMER  -  about to pop...\n")
            v, err := c.Do("RPOP", channel)
            fmt.Printf("CONSUMER  -  Message: %s\n", v)
            if err == nil && v != nil {
                fmt.Printf("\nCONSUMER  -  Successfully Processed Message: %s\n\n", v)
            } else {
                if err != nil {
                    panic("CONSUMER  -  Error: " + err.Error() + "\n")
                }
                fmt.Printf("CONSUMER  -  stop listening\n")
                break
            }
        }
        <- notify
        fmt.Printf("CONSUMER  -  received notification!\n")
    }
}

func publish(channel string, v string) {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil { panic("PUBLISHER  -  Could not connect to Redis") }

    fmt.Printf("PUBLISHER  -  connected - trying to send '%s'\n", v)
    c.Do("LPUSH", channel, v)
    fmt.Printf("PUBLISHER  -  pushed '%s'\n", v)
    c.Do("PUBLISH", channel, v)
    fmt.Printf("PUBLISHER  -  published '%s'\n", v)
}

func sleep_one_second() {
    one_second, err := time.ParseDuration("1s")
    if err != nil { panic("Could not parse duration string '1s'") }
    time.Sleep(one_second)
}

func start_pinging(value []byte) {
    for {
        fmt.Printf("%s...\n", value)
        sleep_one_second()
    }
}

func main() {
    go consume("noise")

    sleep_one_second()

    messages := [4]string{"Here am a message", "Here's another message", "And here's one last one", "Just kidding here's a fourth one lol"}
    
    for _, v := range messages {
        publish("noise", v)
    }

    sleep_one_second()
}