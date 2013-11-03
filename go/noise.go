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
    one_second, err := time.ParseDuration("1s")
    if err != nil { panic("Could not parse duration string '1s'") }
    time.Sleep(one_second)
}

func sleep_forever() {
    for { sleep_one_second() }
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
        p.Publish(v)
    }

    sleep_forever()
}