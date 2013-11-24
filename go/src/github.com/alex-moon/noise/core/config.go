package core

import (
    "code.google.com/p/gcfg"
)

type config struct {
    Redis struct {
        Address string
    }

    Files struct {
        Texts string
    }

    Queues struct {
        Texts string
        Terms string
        Facts string
    }

    Sets struct {
        Mean string
        SD string
        Count string
    }

    SetPrefix struct {
        Correlation string
        CorrelationCount string
        MultiCorrelation string
        MultiCorrelationCount string
        Interest string
    }

    Mutex struct {
        Prefix string
        TTL int
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