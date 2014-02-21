package com.github.alex_moon.noise.core;

import java.lang.reflect.Method;
import java.util.List;

public abstract class Updateable {
	protected List<Updateable> listeners;
	private Method doUpdateMethod;

	public void update(Updateable sender) {
    	try {
			synchronized (this) {
				doUpdate(sender);
				for (Updateable listener : listeners) {
					listener.update(this);
				}
			}
	    } catch (Exception e) {
    		e.printStackTrace();
    	}
	}
	
	public void listen(Updateable obj) {
		if (!listeners.contains(obj)) {
			listeners.add(obj);
		}
	}
	
	protected abstract void doUpdate(Updateable sender);
}
