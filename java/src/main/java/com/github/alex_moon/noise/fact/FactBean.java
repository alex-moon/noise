package com.github.alex_moon.noise.fact;

public class FactBean {
    // simplified version of Fact without all the methods etc.
    private String primaryTerm, x, y;
    private Double mvc;

    public FactBean(Fact fact) {
        primaryTerm = fact.getPrimaryTerm().toString();
        x = fact.getX().toString();
        y = fact.getY().toString();
        mvc = fact.doubleValue();
    }
    
    public String getPrimaryTerm() {
        return primaryTerm;
    }
    
    public String getX() {
        return x;
    }
    
    public String getY() {
        return y;
    }
    
    public Double getMvc() {
        return mvc;
    }
}
