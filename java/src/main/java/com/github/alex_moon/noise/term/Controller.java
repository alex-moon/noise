package com.github.alex_moon.noise.term;

import java.util.HashMap;
import java.util.Map;

import com.github.alex_moon.noise.text.Text;

public class Controller extends Thread {
    private Map<String, Term> terms;

    public void run() {
        terms = new HashMap<String, Term>();
    }

    public Term getTerm(String termString, Text text) {
        Term term = null;
        if (!terms.containsKey(termString)) {
            term = new Term(termString);
            terms.put(termString, term);
        } else {
            term = terms.get(termString);
        }
        text.listen(term);
        return term;
    }
}