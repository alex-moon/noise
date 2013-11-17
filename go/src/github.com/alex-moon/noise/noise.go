package main

import (
    "github.com/alex-moon/noise/core"
    "github.com/alex-moon/noise/texts"
    // "github.com/alex-moon/gostat/stat"
    // we can prove this is imported thus: fmt.Printf("STAT DEBUG %s\n", stat.Binomial_PMF(0.09, 9)(8))
)

type MappedNGram struct {
    key string
    value string
}

func main() {
    texts.TextController()

    select {}
}
