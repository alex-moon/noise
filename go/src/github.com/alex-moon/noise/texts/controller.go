package texts

import (
    "github.com/alex-moon/noise/core"
)

type TextProcessor struct {
    publisher core.Publisher
}

func (p TextProcessor) Process() {
    var text_dir string = core.Config().Files.Texts

    // TODO: make configurable - filesystem as opposed to redis/SQL/RSS/whatever
    iterator := core.NewFileSystemIterator(text_dir)
    for uuid := range iterator.Items() {
        if uuid == nil { break }

        // TODO: ditto
        reader := NewFileSystemReader(text_dir, string(uuid.(string)))
        term_counter := NewTermCounter(reader)
        go term_counter.Run(p.publisher)
    }
}

// TODO I'm not convinced by this naming scheme...
func Texts() {
    consumer := core.NewConsumer(core.Config().Queues.Texts)
    publisher := core.NewPublisher(core.Config().Queues.Terms)
    processor := TextProcessor {publisher}
    go consumer.Consume(processor)
}
