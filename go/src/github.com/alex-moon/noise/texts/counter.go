package texts

import (
    "fmt"
    "regexp"
    "strings"
    "github.com/garyburd/redigo/redis"
    "github.com/alex-moon/noise/core"
    "github.com/alex-moon/noise/terms"
)

type TermCounter struct {
    conn redis.Conn
    text Reader
}

func NewTermCounter(text Reader) TermCounter {
    c, err := redis.Dial("tcp", core.Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("WORD COUNTER %s  -  Could not connect to Redis", text.Uuid()))
    }

    c.Do("DEL", text.Uuid())  // if we're doing this a second time the first time failed - idempotence bitch

    return TermCounter {
        conn: c,
        text: text,
    }
}

func (c TermCounter) Run(publisher core.Publisher) {
    pattern, err := regexp.Compile("[\\W]+")
    if err != nil { panic(err) }

    text_content := strings.ToLower( string( pattern.ReplaceAll(c.text.Bytes(), []byte(" ")) ) )
    words := strings.Split(text_content, " ")

    for _, word := range words {
        if word != "" {
            // fmt.Printf("%s-%s\n", c.text.Uuid(), word)
            c.conn.Do("ZINCRBY", c.text.Uuid(), 1, word)
            c.conn.Do("ZINCRBY", c.text.Uuid(), 1, core.SET_SUM_MEMBER)
            terms.TermLocker().Create(word)
        }
    }

    publisher.Publish(c.text.Uuid())
}
