package terms

import (
    "fmt"
    "github.com/alex-moon/noise/core"
    "github.com/garyburd/redigo/redis"
)

type TextCorrelator struct {
    conn redis.Conn
    iterator core.Iterator
}

func NewTextCorrelator(term_iterator core.Iterator) TextCorrelator {
    c, err := redis.Dial("tcp", core.Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("TEXT CORRELATOR %s  -  Could not connect to Redis", term_iterator))
    }

    return TextCorrelator{c, term_iterator}
}

func (tc TextCorrelator) Run(p core.Publisher) {
    for item := range tc.iterator.Items() {
        if item == nil { break }
        // "member" is the term whose vital stats we'll update at the same time
        member := item.(core.SetMember)

        // TODO designating these as new vars is a waste of memory
        term := member.Term
        count := member.Score
        total := member.SumTotal
        score := count / total
        fmt.Printf("Dig it: term %s appears %s times out of %s for a score of %s\n", term, count, total, score)

        getter := core.NewGetter()
        go func() {
            TermLocker().Lock(term)
            defer TermLocker().Unlock(term)

            // TODO provide defaults for all the following - perhaps some kind of set getter... yeah

            // STEP 1: the count
            // the number of texts this term has appeared in
            old_n :=  getter.Get(core.Config().Sets.Count, term, 0)
            new_n := old_n.(int) + 1

            // STEP 2: the mean
            // M(k) = M(k-1) + (x(k) - M(k-1)) / k
            old_mean := getter.Get(core.Config().Sets.Mean, term, score)
            new_mean := old_mean.(float32) + (score - old_mean.(float32)) / float32(new_n)  // Knuth-Welford

            // STEP 3: the standard deviation
            // S(k) = S(k-1) + (x(k) - M(k-1)) * (x(k) - M(k))
            old_sd := getter.Get(core.Config().Sets.SD, term, 0)
            new_sd := old_sd.(float32) + (score - old_mean.(float32)) * (score - new_mean)
            fmt.Printf("We have a standard deviation: %s\n", new_sd)

            for _, cross_reference_member := range member.CrossReference {
                fmt.Printf("HOLY SHIT - now correlating %s and %s\n", member.Term, cross_reference_member.Term)

                // STEP 4: correlation count
                // bear in mind: everything from here on in is symmetrical
                tc.conn.Do("ZINCRBY", core.Config().SetPrefix.CorrelationCount + term, 1, cross_reference_member.Term)
                tc.conn.Do("ZINCRBY", core.Config().SetPrefix.CorrelationCount + cross_reference_member.Term, 1, term)

                // STEP 5: the moment of truth
                existing_correlation := getter.Get(core.Config().SetPrefix.Correlation + term, cross_reference_member.Term, 0)

                if existing_correlation != nil {
                    cr_sd := getter.Get(core.Config().Sets.SD, cross_reference_member.Term, 0)
                    fmt.Printf("We have a standard deviation: %s\n", cr_sd)
                } else {
                    // all hell breaks loose
                }


                /*
                core.Config().SetPrefix.Correlation + term

                we now have member - here is where we would iterate through all the succeeding set members and do correlation
                if existing correlation:
                1. multiply by product of existing SDs to get existing covariance,
                2. update covariance - http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Covariance (scroll to the very bottom)
                - C(n)(x,y) = C(n-1)(x,y) * (n-1) + (x(n) - M(n)(x)) * (y(n) - M(n-1)(y)) [= same thing but with M(n-1)(x) and M(n)(y)]
                              -----------------------------------------------------------
                                                      n
                3. divide by product of new SDs

                tc.conn.Do("LPUSH", core.Config().Mutex.Prefix + term, "buttlol")
                */
            }
        }()
    }
}