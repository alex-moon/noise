package main

import (
    "fmt"
    // "github.com/fzzy/radix/redis"
    "github.com/garyburd/redigo/redis"
)

type MappedNGram struct {
    key string
    value string
}

func main() {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil { /* whatever */ }

    c.Do("SET", "hello", "world")

    hello, err := redis.String(c.Do("GET", "hello"))
    if err != nil { /* whatever */ }
    fmt.Println(hello)
}
