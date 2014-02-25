package com.github.alex_moon.noise.core;

import java.util.HashMap;
import java.util.Map;

public class Core {
    public static final Integer TEXTS = 0;
    public static final Integer TERMS = 1;
    public static final Integer CORRELATIONS = 2;

    private Map<Integer, Thread> controllers;

    protected Core() {
        controllers = new HashMap<Integer, Thread>();
        controllers.put(TEXTS, new com.github.alex_moon.noise.text.Controller());
        controllers.put(TERMS, new com.github.alex_moon.noise.term.Controller());
        controllers.put(CORRELATIONS, new com.github.alex_moon.noise.correlation.Controller());
    }

    private static Core instance;

    public static Core getInstance() { // We all know the drill.
        if (instance == null) {
            instance = new Core();
        }
        return instance;
    }

    public void run() {
        // Start up all our controllers
        for (Integer type : controllers.keySet()) {
            controllers.get(type).start();
        }
    }

    public static Thread getController(Integer type) {
        return getInstance().controllers.get(type);
    }
    
    public static com.github.alex_moon.noise.text.Controller getTextController() {
        return (com.github.alex_moon.noise.text.Controller) getController(TEXTS);
    }
    
    public static com.github.alex_moon.noise.term.Controller getTermController() {
        return (com.github.alex_moon.noise.term.Controller) getController(TERMS);
    }

    public static com.github.alex_moon.noise.correlation.Controller getCorrelationController() {
        return (com.github.alex_moon.noise.correlation.Controller) getController(CORRELATIONS);
    }
}