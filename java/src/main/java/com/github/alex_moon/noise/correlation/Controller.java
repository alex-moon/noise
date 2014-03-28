package com.github.alex_moon.noise.correlation;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.github.alex_moon.noise.term.Term;

public class Controller extends Thread {
    private List<Correlation> correlations;
    private Map<Term, Map<Term, Correlation>> correlationMap;

    public void run() {
        correlations = new ArrayList<Correlation>();
        correlationMap = new HashMap<Term, Map<Term, Correlation>>();
    }

    // @todo the symmetry problem here has been solved, I'm sure, more
    // efficiently than this by someone else...
    private Correlation addCorrelation(Term a, Term b) {
        Correlation correlation = new Correlation(a, b);
        correlations.add(correlation);

        // first, map from a to b
        if (correlationMap.containsKey(a)) {
            correlationMap.get(a).put(b, correlation);
        } else {
            Map<Term, Correlation> correlationMapRow = new HashMap<Term, Correlation>();
            correlationMapRow.put(b, correlation);
            correlationMap.put(a, correlationMapRow);
        }

        // second, map from b back to a - this ensures a lookup on either will
        // get the same object
        if (correlationMap.containsKey(b)) {
            correlationMap.get(b).put(a, correlation);
        } else {
            Map<Term, Correlation> correlationMapRow = new HashMap<Term, Correlation>();
            correlationMapRow.put(a, correlation);
            correlationMap.put(b, correlationMapRow);
        }

        return correlation;
    }

    public Correlation getCorrelation(Term a, Term b) {
        if (correlationMap.containsKey(a) && correlationMap.get(a).containsKey(b)) {
            return correlationMap.get(a).get(b);
        }
        return addCorrelation(a, b);
    }

    public List<Term> getThirdTerms(Term a, Term b) {
        List<Term> cs = new ArrayList<Term>(a.getCorrelates());
        cs.retainAll(b.getCorrelates());
        return cs;
    }
}
