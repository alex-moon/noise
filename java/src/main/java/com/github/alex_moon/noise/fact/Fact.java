package com.github.alex_moon.noise.fact;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.core.Updateable;
import com.github.alex_moon.noise.term.Term;

public class Fact extends Updateable {
    private Term primaryTerm, x, y;
    
    public Fact(Term primaryTerm, Term x, Term y) {
        this.primaryTerm = primaryTerm;
        this.x = x;
        this.y = y;
        Core.getCorrelationController().getCorrelation(primaryTerm, x).listen(this);
        Core.getCorrelationController().getCorrelation(primaryTerm, y).listen(this);
        Core.getCorrelationController().getCorrelation(x, y).listen(this);
    }

    @Override
    protected void doUpdate(Updateable sender) {
        // see "Formula for R" at http://mtweb.mtsu.edu/stats/regression/level3/multicorrel/multicorrcoef.htm
        // where y is "primaryTerm" and x1 and x2 are x and y in no particular order
    }

}
