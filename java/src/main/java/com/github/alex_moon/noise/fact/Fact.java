package com.github.alex_moon.noise.fact;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.core.Updateable;
import com.github.alex_moon.noise.correlation.Correlation;
import com.github.alex_moon.noise.term.Term;

public class Fact extends Updateable {
    private Term primaryTerm, x, y;
    private Correlation ax, ay, xy;  // where a is primaryTerm
    private Double multivariateCorrelation = 0.0;
    
    public Fact(Term primaryTerm, Term x, Term y) {
        this.primaryTerm = primaryTerm;
        this.x = x;
        this.y = y;
        ax = Core.getCorrelationController().getCorrelation(primaryTerm, x);
        ay = Core.getCorrelationController().getCorrelation(primaryTerm, y);
        xy = Core.getCorrelationController().getCorrelation(x, y);
        ax.listen(this);
        ay.listen(this);
        xy.listen(this);
    }

    @Override
    protected void doUpdate(Updateable sender) {
        // as per http://mtweb.mtsu.edu/stats/regression/level3/multicorrel/multicorrcoef.htm
        // where y is "primaryTerm" and x1 and x2 are "x" and "y" in no particular order
        Double ax2 = ax.doubleValue() * ax.doubleValue();
        Double ay2 = ay.doubleValue() * ay.doubleValue();
        Double xy2 = xy.doubleValue() * xy.doubleValue();
        Double axayxy = ax.doubleValue() * ay.doubleValue() * xy.doubleValue();
        multivariateCorrelation = Math.sqrt(
            ax2 + ay2 - 2 * axayxy
        ) / Math.sqrt(
            1 - xy2
        );

        // in reality these values will need to be tweaked and applied to Interests rather than Facts
        if (multivariateCorrelation > 0.8 && multivariateCorrelation < 0.99) {
            System.out.println("We have a fact! " + x + ":" + y + " -> " + primaryTerm + " @ " + multivariateCorrelation);
        }
    }

    public Double doubleValue() {
        return multivariateCorrelation;
    }
    
    public Term getPrimaryTerm() {
        return primaryTerm;
    }
}
