package terms

import (
    "fmt"
    "github.com/alex-moon/noise/core"
    "github.com/garyburd/redigo/redis"
)

type TextCorrelator struct {
    conn redis.Conn
    cross_reference SetCrossReference
}

func NewTextCorrelator(term_iterator core.Iterator) TextCorrelator {
    c, err := redis.Dial("tcp", core.Config().Redis.Address)
    if err != nil {
        panic(fmt.Sprintf("TEXT CORRELATOR %s  -  Could not connect to Redis", term_iterator))
    }

    cross_reference := NewCrossReference(term_iterator)
    return TextCorrelator{c, cross_reference}
}

func (tc TextCorrelator) OnlineCorrelate(member SetCrossReferenceMember, 
                                         cr_member SetCrossReferenceMember,
                                         n float64) {
    if do_over := recover(); do_over == nil {
        return
    }
    defer tc.OnlineCorrelate(member, cr_member, n)
    getter := tc.cross_reference.getter
    old_correlation := getter.GetFloat(core.Config().SetPrefix.Correlation + member.Term, cr_member.Term, 0.0)
    old_covariance := old_correlation * cr_member.Old.SD * member.Old.SD
    new_covariance := (old_covariance * n + (member.Score - member.New.Mean) *
                       (cr_member.Score - cr_member.Old.Mean)) / n // Pébay
    new_correlation := new_covariance / (cr_member.New.SD * member.New.SD)

    tc.conn.Do("ZADD", core.Config().SetPrefix.Correlation + member.Term, new_correlation, cr_member.Term)
    tc.conn.Do("ZADD", core.Config().SetPrefix.Correlation + cr_member.Term, new_correlation, member.Term)
}

func (tc TextCorrelator) FirstCorrelate(member SetCrossReferenceMember, cr_member SetCrossReferenceMember) {
    if do_over := recover(); do_over == nil {
        return
    }
    defer tc.FirstCorrelate(member, cr_member)
    // Pearson for two observations is always 1, -1 or undefined (if either variable doesn't move)
    // For convenience we're going to assume that if either is equal we've got a floating point error
    // the correlation is in fact 1
    new_correlation := 1.0
    if (member.New.Mean > member.Old.Mean && cr_member.New.Mean < cr_member.Old.Mean) ||
       (member.New.Mean < member.Old.Mean && cr_member.New.Mean > cr_member.Old.Mean) {
        new_correlation = -1.0
    }
    tc.conn.Do("ZADD", core.Config().SetPrefix.Correlation + member.Term, new_correlation, cr_member.Term)
    tc.conn.Do("ZADD", core.Config().SetPrefix.Correlation + cr_member.Term, new_correlation, member.Term)
}

func (tc TextCorrelator) Correlate(member SetCrossReferenceMember) {
    getter := tc.cross_reference.getter
    for _, cr_member := range member.CrossReference {
        correlation_count := getter.GetInt(core.Config().SetPrefix.CorrelationCount + member.Term, cr_member.Term, 0)

        if correlation_count > 1 {
            go tc.OnlineCorrelate(member, cr_member, float64(correlation_count))
        } else if correlation_count > 0 {
            go tc.FirstCorrelate(member, cr_member)
        }

        // STEP 5: correlation count
        tc.conn.Do("ZINCRBY", core.Config().SetPrefix.CorrelationCount + member.Term, 1, cr_member.Term)
        tc.conn.Do("ZINCRBY", core.Config().SetPrefix.CorrelationCount + cr_member.Term, 1, member.Term)
    }
}

func (tc TextCorrelator) Run(p core.Publisher) {
    for item := range tc.cross_reference.Items() {
        if item == nil { break }
        go tc.Correlate(item.(SetCrossReferenceMember))
    }
}