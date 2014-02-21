package com.github.alex_moon.noise.core;

public class Core {
	protected Core() {}  // Yes, this class is a singleton. Pretty original.
	private static Core instance;
	
	private com.github.alex_moon.noise.text.Controller textController;
	private com.github.alex_moon.noise.term	.Controller termController;
	private com.github.alex_moon.noise.correlation.Controller correlationController;
	
	public static Core getInstance() {  // We all know the drill.
		if (instance == null) {
			instance = new Core();
		}
		return instance;
	}
	
	public void run() {
		// Start up all our controllers
		(textController = new com.github.alex_moon.noise.text.Controller()).start();
        (termController = new com.github.alex_moon.noise.term.Controller()).start();
        (correlationController = new com.github.alex_moon.noise.correlation.Controller()).start();
	}
	
	public static com.github.alex_moon.noise.text.Controller getTextController() { return getInstance().textController; }	
	public static com.github.alex_moon.noise.term.Controller getTermController() { return getInstance().termController; }	
	public static com.github.alex_moon.noise.correlation.Controller getCorrelationController() { return getInstance().correlationController; }
}
