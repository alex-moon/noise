package main

import (
    "fmt"
    "regexp"
    "strings"
    "github.com/garyburd/redigo/redis"
)

type TermCounter struct {
    reader Reader
    conn redis.Conn
}

func NewTermCounter(reader Reader) TermCounter {
    conn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("WORD COUNTER %s  -  Could not connect to Redis", reader.text_id))
    }
    return TermCounter{reader, conn}
}

func (c TermCounter) Run() {
    pattern, err := regexp.Compile("[\\W]+")
    if err != nil { panic(err) }

    text_content_slug := string(pattern.ReplaceAll(c.reader.ReadAll(), []byte(" ")))
    text_content_slug = strings.ToLower(text_content_slug)

    words := strings.Split(text_content_slug, " ")

    for _, word := range words {
        fmt.Printf("%s-%s\n", c.reader.text_id, word)
        c.conn.Do("ZINCRBY", c.reader.text_id, 1, word)
    }
}