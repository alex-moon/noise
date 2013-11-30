package terms

import (
    "github.com/alex-moon/noise/core"
)

type VitalStats struct {
    N int
    Mean float64
    SD float64
}

type SetCrossReferenceMember struct {
    Term string
    Score float64
    Old VitalStats
    New VitalStats
    CrossReference []SetCrossReferenceMember
}

type SetCrossReference struct {
    getter core.RedisScoreGetter
    iterator core.Iterator
}

func NewCrossReference(iterator core.Iterator) SetCrossReference {
    getter := core.NewGetter()
    return SetCrossReference{getter, iterator}
}

func (cr SetCrossReference) Items() chan core.Item {
    var items_chan chan core.Item = make(chan core.Item)

    go func() {
        var items []SetCrossReferenceMember = []SetCrossReferenceMember{}

        for item := range cr.iterator.Items() {
            if item == nil { break }
            member := item.(core.SetMember)
            score := member.Score / member.SumTotal

            // STEP 1: the count
            // the number of texts this item.Term has appeared in
            old_n :=  cr.getter.GetInt(core.Config().Sets.Count, member.Term, 0)
            new_n := old_n + 1

            // STEP 2: the mean
            // M(k) = M(k-1) + (x(k) - M(k-1)) / k
            old_mean := cr.getter.GetFloat(core.Config().Sets.Mean, member.Term, score)
            new_mean := old_mean + (score - old_mean) / float64(new_n)  // Knuth-Welford

            // STEP 3: the standard deviation
            // S(k) = S(k-1) + (x(k) - M(k-1)) * (x(k) - M(k))
            old_sd := cr.getter.GetFloat(core.Config().Sets.SD, member.Term, 0.0)
            new_sd := old_sd + (score - old_mean) * (score - new_mean)  // Knuth-Welford

            items = append(items, SetCrossReferenceMember{
                member.Term,
                score,
                VitalStats {
                    old_n,
                    old_mean,
                    old_sd,
                },
                VitalStats {
                    new_n,
                    new_mean,
                    new_sd,
                },
                []SetCrossReferenceMember{},
            })
        }

        for i := range items {
            items[i].CrossReference = items[i+1:]
            items_chan <- items[i]
        }
    } ()

    return items_chan
}
