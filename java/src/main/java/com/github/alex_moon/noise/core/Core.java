package com.github.alex_moon.noise.core;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Scanner;

public class Core {
    public static final String textsDir = "/home/moona/work/noise/text";
    private final String stopwordsFile = "/home/moona/work/noise/stopwords.txt";
    public static List<String> stopwords;
    
    public static final Integer TEXTS = 0;
    public static final Integer TERMS = 1;
    public static final Integer CORRELATIONS = 2;
    public static final Integer FACTS = 3;

    private Map<Integer, Thread> controllers;

    protected Core()  { // Something something singleton yada yada yada...
        // first let's load those stopwords
        try {
            Scanner s = new Scanner(new File(stopwordsFile));
            stopwords = new ArrayList<String>();
            while (s.hasNext()){
                stopwords.add(s.next());
            }
            s.close();
        } catch(FileNotFoundException e) {
            System.out.println(stopwordsFile + " couldn't be found. Are you sure it's there?");
            System.exit(1);
        }
        
        // finally flip all switches to "on" position...        
        controllers = new HashMap<Integer, Thread>();
        controllers.put(TEXTS, new com.github.alex_moon.noise.text.Controller());
        controllers.put(TERMS, new com.github.alex_moon.noise.term.Controller());
        controllers.put(CORRELATIONS, new com.github.alex_moon.noise.correlation.Controller());
        controllers.put(FACTS, new com.github.alex_moon.noise.fact.Controller());
    }

    private static Core instance;

    public static Core getInstance()  { // We all know the drill.
        if (instance == null) {
            instance = new Core();
            instance.run();  // so far I can't think of a use case where we don't need it running...
        }
        return instance;
    }

    public void run() {
        // Start up all our controllers
        for (Integer type : controllers.keySet()) {
            controllers.get(type).start();
        }
    }
    
    public Boolean isStopword(String term) {
        return stopwords.contains(term);
    }

    public static Thread getController(Integer type)  {
        return getInstance().controllers.get(type);
    }

    public static com.github.alex_moon.noise.text.Controller getTextController()  {
        return (com.github.alex_moon.noise.text.Controller) getController(TEXTS);
    }

    public static com.github.alex_moon.noise.term.Controller getTermController()  {
        return (com.github.alex_moon.noise.term.Controller) getController(TERMS);
    }

    public static com.github.alex_moon.noise.correlation.Controller getCorrelationController()  {
        return (com.github.alex_moon.noise.correlation.Controller) getController(CORRELATIONS);
    }

    public static com.github.alex_moon.noise.fact.Controller getFactController()  {
        return (com.github.alex_moon.noise.fact.Controller) getController(FACTS);
    }
}