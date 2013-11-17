package core

type Processor interface {
    NewIterator() chan interface {}
    NewWorker() func(interface {})
}