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
    fmt.Printf("New Term Counter with text %s\n", text)

    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("WORD COUNTER %s  -  Could not connect to Redis", text.Uuid()))
    }

    fmt.Printf("About to return Term Counter\n")
    return TermCounter{c, text}
}

func (c TermCounter) Run() {
    pattern, err := regexp.Compile("[\\W]+")
    if err != nil { panic(err) }

    fmt.Printf("We have a pattern %s\n", pattern)

    text_content_slug := string(pattern.ReplaceAll(c.text.Bytes(), []byte(" ")))

    fmt.Printf("We have a pattern replace %s\n", text_content_slug)
    
    text_content_slug = strings.ToLower(text_content_slug)

    fmt.Printf("We have a slug %s\n", text_content_slug)

    words := strings.Split(text_content_slug, " ")

    for _, word := range words {
        fmt.Printf("%s-%s\n", c.text.Uuid(), word)
        c.conn.Do("ZINCRBY", c.text.Uuid(), 1, word)
    }
}