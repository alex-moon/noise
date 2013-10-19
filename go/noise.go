package main

import (
    "fmt"
    "time"
)

type MappedNGram struct {
    key string
    value string
}

func sleep_one_second() {
    one_second, err := time.ParseDuration("10s")
    if err != nil { panic("Could not parse duration string '10s'") }
    time.Sleep(one_second)
}

func start_pinging(value []byte) {
    for {
        fmt.Printf("%s...\n", value)
        sleep_one_second()
    }
}

func main() {
    c := NewConsumer("noise")
    p := NewPublisher("noise")

    go c.Consume()

    messages := [4]string{"Here am a message", "Here's another message", "And here's one last one", "Just kidding here's a fourth one lol"}
    
    for _, v := range messages {
        p.Publish("noise", v)
    }

    sleep_one_second()
}