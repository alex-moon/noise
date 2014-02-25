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
        // only update the correlation if both terms have been hit by the same text
        if (a.getLastText() == b.getLastText()) {
            Integer n = getN();
            if (n > 2) {
                Double aDelta = a.getProportion() - a.getOldMean();
                Double bDelta = b.getProportion() - b.getOldMean();
                Double aMean = a.getOldMean() + aDelta / (n + 1);
                Double bMean = b.getOldMean() + bDelta / (n + 1);
                System.out.println("New mean for " + a + " should be " + aMean + " is actually " + a.getMean());
                System.out.println("New mean for " + b + " should be " + bMean + " is actually " + b.getMean());
                Double covariance = n * coefficient * a.getOldSd() * b.getOldSd() +
                                    n * aDelta * bDelta / (n + 1);
                coefficient = covariance / (a.getSd() * b.getSd()); // should these be oldSd()?
                if (a.toString() == b.toString()) System.out.println("Term '" + a + "' (sd=" + a.getSd() + ") and '" + b + "' (sd="+ b.getSd() +"): n=" + getN() + " r=" + coefficient);
            } else {
                coefficient = 1.0;
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
}