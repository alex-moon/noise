package main

import (
    "code.google.com/p/gcfg"
)

type config struct {
    Text struct {
        Dir string
    }
}

var config_instance *config = nil

func Config() *config {
    if config_instance == nil {
        config_instance = new(config);
        err := gcfg.ReadFileInto(config_instance, "noise.ini")
        if err != nil { panic(err) }
    }
    return config_instance
}