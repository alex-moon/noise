package texts

import (
    "github.com/alex-moon/noise/core"
)

type TextProcessor struct {
    publisher core.Publisher
}

func (p TextProcessor) Process() {
    reader := NewReader()
    for text := range reader.texts {
        if text == nil { break }
        term_counter := NewTermCounter(text)
        go term_counter.Run(p.publisher)
    }
}

// TODO noun-named method should return something
func Texts() {
    consumer := core.NewConsumer(core.Config().Queues.Texts)
    publisher := core.NewPublisher(core.Config().Queues.Terms)
    processor := TextProcessor {publisher}
    go consumer.Consume(processor)
}