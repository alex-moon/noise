package com.github.alex_moon.noise.term;

import java.util.ArrayList;
import java.util.List;

import com.github.alex_moon.noise.core.Updateable;
import com.github.alex_moon.noise.text.Text;

public class Term extends Updateable {
    private String termString;
    private Double mean = 0.0;
    private Double standardDeviation = 0.0;
    private List<Text> texts = new ArrayList<Text>();

    private Text lastText;
    private Double lastMean, lastStandardDeviation, lastProportion;

    public Term(String termString) {
        this.termString = termString;
    }

    /**
     * @param lastProportion
     *            = the lastProportion of some text that is the term = term word
     *            count / total word count
     */
    protected void doUpdate(Updateable text) {
        lastText = (Text) text;
        lastProportion = lastText.getProportion(this.termString);
        if (texts.contains(lastText)) {
            // @todo something?
            System.out.println("WARNING: text " + lastText.getUuid()
                    + " reregistering for term '" + termString + "'");
        } else {
            texts.add(lastText);
            Integer count = texts.size();
            lastMean = mean;
            lastStandardDeviation = standardDeviation;
            mean = mean + (lastProportion - mean) / count;
            Double sumOfSquaredDifferences = (standardDeviation * standardDeviation)
                    * (count - 1)
                    + ((lastProportion - lastMean) * (lastProportion - mean));
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

    public Double getProportion() {
        return lastProportion;
    }
}