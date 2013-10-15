package main

import (
    "fmt"
    "github.com/garyburd/redigo/redis"
)

type Db struct {
    conn redis.Conn
}

func NewDb() Db {
    c, err := redis.Dial("tcp", ":6379")
    if err != nil {
        panic(fmt.Sprintf("DB  -  Could not connect to Redis"))
    }

    return Db{
        conn: c,
    }
}

func (db Db) Get(key string, default_value string) {
    value, err := s.conn.Do("GET", key)
    if err == nil {
        if value == nil {
            return default_value
        } else {
            return value
        }
    } else {
        panic(fmt.Sprintf("DB  -  Get error {%s '%s'}: %s", key, default_value, err.Error()))
    }
}

func (db Db) Set(key string, value string) {
    _, err := s.conn.Do("SET", key, value)
    if err != nil {
        panic(fmt.Sprintf("DB  -  Set error {%s '%s'}: %s", key, value, err.Error()))
    }
}