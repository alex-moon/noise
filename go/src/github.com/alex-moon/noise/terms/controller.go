package terms

import (
    "github.com/alex-moon/noise/core"
)

type TermProcessor struct {
    publisher core.Publisher
}

func (p TermProcessor) Process() {
    iterator := core.NewQueueIterator(core.Config().Queues.Terms)
    for uuid := range iterator.items {
        if uuid == nil { break }
        term_iterator := core.NewSetIterator(uuid)
        term_correlator := NewTermCorrelator(term_iterator)  // TODO: ZRANGE over terms and generate SDs, means, etc. - this is a bit more complex
        go term_correlator.Run(p.publisher)
    }
}

func Terms() {
    consumer := core.NewConsumer(core.Config().Queues.Terms)
    publisher := core.NewPublisher(core.Config().Queues.Facts)
    processor := TermProcessor {publisher}
    go consumer.Consume(processor)
}