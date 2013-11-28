package terms

import (
    "github.com/alex-moon/noise/core"
)

type CrossReference core.Iterator

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
    items []SetCrossReferenceMember
}

func NewCrossReference(iterator core.Iterator) CrossReference {
    var items []SetCrossReferenceMember = []SetCrossReferenceMember{}
    
    // TODO: use recover() to check if we're in the middle of a panic?
    // TermLocker().Lock(member.Term)
    // defer TermLocker().Unlock(member.Term)
    getter := core.NewGetter()

    for item := range iterator.Items() {
        if item == nil { break }
        member := item.(core.SetMember)
        score := member.Score / member.SumTotal

        // STEP 1: the count
        // the number of texts this item.Term has appeared in
        old_n :=  getter.GetInt(core.Config().Sets.Count, member.Term, 0)
        new_n := old_n + 1

        // STEP 2: the mean
        // M(k) = M(k-1) + (x(k) - M(k-1)) / k
        old_mean := getter.GetFloat(core.Config().Sets.Mean, member.Term, score)
        new_mean := old_mean + (score - old_mean) / float64(new_n)  // Knuth-Welford

        // STEP 3: the standard deviation
        // S(k) = S(k-1) + (x(k) - M(k-1)) * (x(k) - M(k))
        old_sd := getter.GetFloat(core.Config().Sets.SD, member.Term, 0.0)
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
    }

    return SetCrossReference{items}
}

func (cr SetCrossReference) Items() chan core.Item {
    items := make(chan core.Item)
    go func() {
        for _, item := range cr.items {
            items <- item
        }
    } ()
    return items
}
