package com.github.alex_moon.noise.text;

public class Controller extends Thread {
    private Retriever retriever;

    public Controller() {
        retriever = new Retriever();
    }

    public void run() {
        retriever.start();
    }
}