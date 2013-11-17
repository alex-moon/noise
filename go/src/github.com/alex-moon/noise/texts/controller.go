package texts

import (
    "github.com/alex-moon/noise/core"
)

type TextProcessor struct {
    publisher core.Publisher
}

func (p TextProcessor) NewWorker() func(TextReader) {
    return func(text TextReader) {
        counter := NewTermCounter(text)
        counter.Run(p.publisher)
    }
}

func (p TextProcessor) NewIterator() chan TextReader {
    reader := NewReader()
    return reader.texts
}

// TODO noun-named method should return something
func TextsController() {
    consumer := core.NewConsumer(core.Config().Queues.Texts)
    publisher := core.NewPublisher(core.Config().Queues.Terms)
    processor := TextProcessor {publisher}
    go consumer.Consume(processor)
}