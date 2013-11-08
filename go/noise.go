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

    p := NewPublisher("noise")  // on the "noise" list...
    p.Publish("aee02730-47e7-11e3-8f96-0800200c9a66")  // ...the text with this ID is available for processing

    for {}
}
