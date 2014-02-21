package com.github.alex_moon.noise.term;

import java.util.HashMap;
import java.util.Map;

public class Controller extends Thread {
	private Map<String, Term> terms;
    public void run() {
    	terms = new HashMap<String, Term>();
        // @todo stuff re terms
    }
    
    public Term getTerm(String termString) {
    	if (!terms.containsKey(termString)) {
    		terms.put(termString, new Term(termString));
    	}
    	return terms.get(termString);
    }
}