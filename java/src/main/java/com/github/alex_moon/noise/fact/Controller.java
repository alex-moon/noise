package com.github.alex_moon.noise.fact;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.github.alex_moon.noise.term.Term;

public class Controller extends Thread {
    private List<Fact> facts;
    private Map<Term, Map<Term, Map<Term, Fact>>> factMap;
    
    private Map<Term, List<Fact>> factsForPrimaryTermCache;

    public void run() {
        facts = new ArrayList<Fact>();
        factMap = new HashMap<Term, Map<Term, Map<Term, Fact>>>();
    }
    
    public List<Fact> getFactsForPrimaryTerm(Term primaryTerm) {
        List<Fact> factsForPrimaryTerm = new ArrayList<Fact>();
        if (factMap.containsKey(primaryTerm)) {
            for (Map<Term, Fact> factValues: factMap.get(primaryTerm).values()) {
                factsForPrimaryTerm.addAll(factValues.values());
            }
        }
        // factsForPrimaryTermCache.set(primaryTerm, factsForPrimaryTerm);  // @todo expire on addFact(), instantiate on run(), etc.
        return factsForPrimaryTerm;
    }

    public Fact getFact(Term primaryTerm, Term x, Term y) {
        if (
            factMap.containsKey(primaryTerm) && 
            factMap.get(primaryTerm).containsKey(x) &&
            factMap.get(primaryTerm).get(x).containsKey(y)
        ) {
            return factMap.get(primaryTerm).get(x).get(y);
        }
        return addFact(primaryTerm, x, y);
    }
    
    private Fact addFact(Term primaryTerm, Term x, Term y) {
        if (!factMap.containsKey(primaryTerm)) {
            factMap.put(primaryTerm, new HashMap<Term, Map<Term, Fact>>());
        }
        Map<Term, Map<Term, Fact>> primaryTermFactMap = factMap.get(primaryTerm);
        
        if (!primaryTermFactMap.containsKey(x)) {
            primaryTermFactMap.put(x, new HashMap<Term, Fact>());
        }
        if (!primaryTermFactMap.containsKey(y)) {
            primaryTermFactMap.put(y, new HashMap<Term, Fact>());
        }
        
        Fact fact = null;
        if (primaryTermFactMap.get(x).containsKey(y)) {
            fact = primaryTermFactMap.get(x).get(y);
        } else {
            fact = new Fact(primaryTerm, x, y);
            primaryTermFactMap.get(x).put(y, fact);
            primaryTermFactMap.get(y).put(x, fact);
        }
        return fact;
    }
}