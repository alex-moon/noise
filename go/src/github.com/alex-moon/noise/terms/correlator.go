package terms

import (
    "fmt"
    "github.com/alex-moon/noise/core"
    "github.com/garyburd/redigo/redis"
)

type TextCorrelator struct {
    conn redis.Conn
    cross_reference CrossReference
}

func NewTextCorrelator(term_iterator core.Iterator) TextCorrelator {
    c, err := redis.Dial("tcp", core.Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("TEXT CORRELATOR %s  -  Could not connect to Redis", term_iterator))
    }

    cross_reference := NewCrossReference(term_iterator)
    return TextCorrelator{c, cross_reference}
}

func (tc TextCorrelator) Run(p core.Publisher) {
    getter := core.NewGetter()
    for item := range tc.cross_reference.Items() {
        if item == nil { break }
        // "member" is the term whose vital stats we'll update at the same time
        member := item.(SetCrossReferenceMember)

        go func() {
            for _, cr_member := range member.CrossReference {
                fmt.Printf("Now correlating %s and %s\n", member.Term, cr_member.Term)
                // bear in mind: everything from here on in is symmetrical

                // STEP 4: the moment of truth
                correlation_count := getter.Get(core.Config().SetPrefix.CorrelationCount + member.Term, cr_member.Term, float32(0.0))

                if correlation_count.(float32) > 0.0 {
                    old_correlation := getter.Get(core.Config().SetPrefix.Correlation + member.Term, cr_member.Term, nil)
                    if old_correlation == nil {
                        panic(fmt.Sprintf("TERM CORRELATOR  -  missing correlation for correlation count > 0 - %s and %s have been correlated %d times\n", member.Term, cr_member.Term, correlation_count))
                    }
                    old_covariance := old_correlation.(float32) * cr_member.Old.SD * member.Old.SD
                    new_covariance := (old_covariance * correlation_count.(float32) + (member.Score - member.New.Mean) * (cr_member.Score - cr_member.Old.Mean)) / correlation_count.(float32) // Pébay
                    new_correlation := new_covariance / (cr_member.New.SD * member.New.SD)

                    tc.conn.Do("ZADD", core.Config().SetPrefix.Correlation + member.Term, new_correlation, cr_member.Term)
                    tc.conn.Do("ZADD", core.Config().SetPrefix.Correlation + cr_member.Term, new_correlation, member.Term)
                } else {
                    // TODO
                    // Pearson for two observations is always 1 or -1 - confirm this works for online covariance above
                    // if it does we can optimise this else case - if they move in the same direction it's 1, otherwise it's -1
                    // otherwise I suspect we're going to have to store the first two scores for every term :(
                }

                // STEP 5: correlation count
                tc.conn.Do("ZINCRBY", core.Config().SetPrefix.CorrelationCount + member.Term, 1, cr_member.Term)
                tc.conn.Do("ZINCRBY", core.Config().SetPrefix.CorrelationCount + cr_member.Term, 1, member.Term)
            }
        }()
    }
}