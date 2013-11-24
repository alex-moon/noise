package terms

import (
    "github.com/alex-moon/noise/core"
)

type TermProcessor struct {
    publisher core.Publisher
}

func (p TermProcessor) Process() {
    iterator := core.NewQueueIterator(core.Config().Queues.Terms)
    for uuid := range iterator.Items() {
        if uuid == nil { break }
        term_iterator := core.NewSetIterator(string(uuid.([]byte)), core.SET_RANK_ITERATOR)
        term_correlator := NewTextCorrelator(term_iterator)
        go term_correlator.Run(p.publisher)
    }
}

func Terms() {
    consumer := core.NewConsumer(core.Config().Queues.Terms)
    publisher := core.NewPublisher(core.Config().Queues.Facts)
    processor := TermProcessor {publisher}
    go consumer.Consume(processor)
}