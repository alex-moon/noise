package com.github.alex_moon.noise.term;

import java.util.ArrayList;
import java.util.List;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.core.Updateable;
import com.github.alex_moon.noise.correlation.Correlation;
import com.github.alex_moon.noise.text.Text;

public class Term extends Updateable {
    private String termString;
    private Double mean = 0.0;
    private Double standardDeviation = 0.0;
    private List<Text> texts = new ArrayList<Text>();
    private List<Term> correlates = new ArrayList<Term>();

    private Text lastText;
    private Double lastMean, lastStandardDeviation, lastProportion;

    public Term(String termString) {
        this.termString = termString;
    }

    /**
     * @param lastProportion
     *            the amount of the last encountered text this term accounts for by itself 
     *            i.e. no. occurrences of this word / no. words in the whole text
     */
    @Override
    protected void doUpdate(Updateable text) {
        lastText = (Text) text;
        lastProportion = lastText.getProportion(this.termString);
        if (texts.contains(lastText)) {
            // @todo something?
            System.out.println("WARNING: text " + lastText.getUuid() + " reregistering for term '" + termString + "'");
        } else {
            texts.add(lastText);

            // first handle correlations
            for (Term correlate : lastText.getCorrelates(this)) {
                if (this == correlate) {
                    continue;
                }
                
                // keep a list of correlates for the facts module
                if (!correlates.contains(correlate)) {
                    correlates.add(correlate);
                }
                Core.getCorrelationController().getCorrelation(this, correlate);
            }

            // now do mean and standard deviation
            Integer count = texts.size();
            lastMean = mean;
            lastStandardDeviation = standardDeviation;
            mean = mean + (lastProportion - mean) / count;
            Double sumOfSquaredDifferences = (standardDeviation * standardDeviation) * (count - 1) + ((lastProportion - lastMean) * (lastProportion - mean));
            standardDeviation = Math.sqrt(sumOfSquaredDifferences / count);
        }
    }

    public Integer count() {
        return texts.size();
    }

    public Double getSd() {
        return standardDeviation;
    }

    public Double getMean() {
        return mean;
    }

    public Double getOldMean() {
        return lastMean;
    }

    public Double getOldSd() {
        return lastStandardDeviation;
    }

    public List<Text> getTexts() {
        return texts;
    }
    
    public List<Term> getCorrelates() {
        return correlates;
    }

    public Text getLastText() {
        return lastText;
    }

    public Double getProportion() {
        return lastProportion;
    }

    public String toString() {
        return termString;
    }
}