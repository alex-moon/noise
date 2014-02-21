package com.github.alex_moon.noise.text;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

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
        	count ++; total ++;
        	counts.put(termString, count);
    	}

    	for (String termString : counts.keySet()) {
    		Term term = new Term(termString);  // @todo getTerm() from term Controller
    		proportions.put(termString, counts.get(termString).doubleValue() / total);
    		listen(term);
    	}
    }
    
    public Double getProportion(String termString) {
    	return proportions.get(termString);
    }

    public String[] asWordList() {
        String result = stringValue.replaceAll("[^a-zA-Z0-9]+", " ");
        result = result.toLowerCase();
        return result.split(" ");
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