package terms

import (
    "fmt"
    "github.com/alex-moon/noise/core"
)

type TermCorrelator struct {
    iterator core.Iterator
}

func NewTermCorrelator(term_iterator core.Iterator) TermCorrelator {
    return TermCorrelator{term_iterator}
}

func (tc TermCorrelator) Run(p core.Publisher) {
    for item := range tc.iterator.Items() {
        member := item.(core.SetMember)
        term := member.Term
        score := member.Score
        fmt.Printf("Dig it: term %s score %s\n", term, score)
    }
}