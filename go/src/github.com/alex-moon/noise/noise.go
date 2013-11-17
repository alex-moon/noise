package main

import (
    // "github.com/alex-moon/gostat/stat"
    // we can prove this is imported thus: fmt.Printf("STAT DEBUG %s\n", stat.Binomial_PMF(0.09, 9)(8))
)

type MappedNGram struct {
    key string
    value string
}

func main() {
    c := NewConsumer("noise")
    go c.Consume()

    p := NewPublisher("noise")
    p.Publish("1")  // 1 new text is available for processing

    select {}
}
