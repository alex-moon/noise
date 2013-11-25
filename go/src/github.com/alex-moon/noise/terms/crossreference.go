package terms

import (
    "github.com/alex-moon/noise/core"
)

type CrossReference core.Iterator

type VitalStats struct {
    N int
    Mean float32
    SD float32
}

type SetCrossReferenceMember struct {
    Term string
    Score float32
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
        old_n :=  getter.Get(core.Config().Sets.Count, member.Term, 0)
        new_n := old_n.(int) + 1

        // STEP 2: the mean
        // M(k) = M(k-1) + (x(k) - M(k-1)) / k
        old_mean := getter.Get(core.Config().Sets.Mean, member.Term, score)
        new_mean := old_mean.(float32) + (score - old_mean.(float32)) / float32(new_n)  // Knuth-Welford

        // STEP 3: the standard deviation
        // S(k) = S(k-1) + (x(k) - M(k-1)) * (x(k) - M(k))
        old_sd := getter.Get(core.Config().Sets.SD, member.Term, float32(0.0))
        new_sd := old_sd.(float32) + (score - old_mean.(float32)) * (score - new_mean)  // Knuth-Welford

        items = append(items, SetCrossReferenceMember{
            member.Term,
            score,
            VitalStats {
                old_n.(int),
                old_mean.(float32),
                old_sd.(float32),
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
