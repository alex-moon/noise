package com.github.alex_moon.noise.correlation;

import com.github.alex_moon.noise.core.Updateable;
import com.github.alex_moon.noise.term.Term;
import com.github.alex_moon.noise.text.Text;

public class Correlation extends Updateable {
    private Term a, b;
    private Double coefficient = 0.0;
    private Integer n;

    public Correlation(Term a, Term b) {
        this.a = a;
        this.b = b;
        a.listen(this);
        b.listen(this);
    }
    
    public void doUpdate(Updateable sender) {
    	Double covariance = coefficient * a.getOldSd() * b.getOldSd();
    	covariance = covariance * getN() + a.getProportion() - a.getMean();
    	covariance = covariance * (b.getProportion() - b.getOldMean()) / getN();
    	coefficient = covariance / (a.getSd() * b.getSd());
    }

    private Integer getN() {
    	if (n != null) {
    		return n;
    	}
    	Integer count = 0;
    	for(Text text : a.getTexts()) {
    		if (b.getTexts().contains(text)) {
    			count ++;
    		}
    	}
    	return count;
    }
}