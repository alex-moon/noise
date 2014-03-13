package com.github.alex_moon.noise.correlation;

import com.github.alex_moon.noise.core.Core;
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

    @Override
    public void doUpdate(Updateable sender) {
        // only update the correlation if both terms have been hit by the same text
        if (a.getLastText() == b.getLastText()) {
            Integer n = getN();
            if (n > 2) {
                Double oldCoefficient = coefficient;
                Double oldCovariance = oldCoefficient * a.getOldSd() * b.getOldSd();
                Double aDelta = a.getProportion() - a.getOldMean();
                Double bDelta = b.getProportion() - b.getOldMean();
                Double newCovariance = oldCovariance + (n-1) * aDelta * bDelta / n;
                coefficient = newCovariance / (a.getSd() * b.getSd());
            } else {
                coefficient = 1.0;
            }

            // now we want to create/update all the relevant Facts, get 'em listen()ing
            for (Term c : Core.getCorrelationController().getThirdTerms(a, b)) {
                // a "third term" c is any term correlated to two given correlated terms
                // we use these to create/update facts - for any three correlated terms
                // there are three facts, one each for each "primary term" (see Fact class)
                listen(Core.getFactController().getFact(a, b, c));
                listen(Core.getFactController().getFact(b, a, c));
                listen(Core.getFactController().getFact(c, a, b));
            }
        }
    }

    private Integer getN() {
        if (n != null) {
            return n;
        }
        Integer count = 0;
        for (Text text : a.getTexts()) {
            if (b.getTexts().contains(text)) {
                count++;
            }
        }
        return count;
    }
    
    public Double doubleValue() {
        return coefficient;
    }
}