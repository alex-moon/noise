package com.github.alex_moon.noise.text;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.core.Updateable;
import com.github.alex_moon.noise.term.Term;

public class Text extends Updateable {
    private String stringValue;
    private UUID uuid;

    private Map<String, Double> proportions;

    public Text(String initialStringValue) {
        stringValue = initialStringValue;
        uuid = UUID.randomUUID();
        proportions = new HashMap<String, Double>();
    }

    public Text(String initialStringValue, String uuid) {
        stringValue = initialStringValue;
        try {
            this.uuid = UUID.fromString(uuid);
        } catch (IllegalArgumentException e) {
            this.uuid = UUID.randomUUID();
        }
        proportions = new HashMap<String, Double>();
    }

    public void doUpdate(Updateable sender) {
        Map<String, Integer> counts = new HashMap<String, Integer>();
        Integer total = 0;

        for (String termString : asWordList()) {
            Integer count = 0;
            if (counts.containsKey(termString)) {
                count += counts.get(termString);
            }
            count++;
            total++;
            counts.put(termString, count);
        }

        for (String termString : counts.keySet()) {
            Core.getTermController().getTerm(termString, this);
            proportions.put(termString, counts.get(termString).doubleValue() / total);
        }
    }

    public List<Term> getCorrelates(Term term) {
        List<Term> correlates = new ArrayList<Term>();
        for (String termString : proportions.keySet()) {
            Term correlate = Core.getTermController().getTerm(termString, this);
            correlates.add(correlate);
        }
        return correlates;
    }

    public Double getProportion(String termString) {
        return proportions.get(termString);
    }

    public String[] asWordList() {
        String result = stringValue.replaceAll("[^a-zA-Z0-9]+", " ");
        result = result.toLowerCase();
        String[] terms = result.split(" ");
        List<String> significantTerms = new ArrayList<String>();
        for (int i = 0; i < terms.length; i++) {
            if (! Core.getInstance().isStopword(terms[i])
             && ! "".equals(terms[i])
             && ! (terms[i].length() < 3)) {
                significantTerms.add(terms[i]);
            }
        }
        return significantTerms.toArray(new String[significantTerms.size()]);
    }

    public String toString() {
        return stringValue;
    }

    public UUID getUuid() {
        return uuid;
    }

    public String getUuidAsString() {
        return uuid.toString();
    }
}