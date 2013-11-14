package main

import (
    "fmt"
    "regexp"
    "strings"
    "github.com/garyburd/redigo/redis"
)

type TermCounter struct {
    conn redis.Conn
    text TextReader
}

func NewTermCounter(text TextReader) TermCounter {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("WORD COUNTER %s  -  Could not connect to Redis", text.Uuid()))
    }
    return TermCounter{c, text}
}

func (c TermCounter) Run() {
    pattern, err := regexp.Compile("[\\W]+")
    if err != nil { panic(err) }

    text_content := strings.ToLower( string( pattern.ReplaceAll(c.text.Bytes(), []byte(" ")) ) )
    words := strings.Split(text_content, " ")

    for _, word := range words {
        if word != "" {
            fmt.Printf("%s-%s\n", c.text.Uuid(), word)
            c.conn.Do("ZINCRBY", c.text.Uuid(), 1, word)
        }
    }
}