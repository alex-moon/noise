package texts

import (
    "github.com/alex-moon/noise/core"
)

func TextsController() {
    consumer := core.NewConsumer(core.Config().Queues.Texts)
    publisher := core.NewPublisher(core.Config().Queues.Terms)
    go consumer.Consume(NewIterator, NewWorker(publisher))
}