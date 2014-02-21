package com.github.alex_moon.noise.core;

import java.lang.reflect.Method;
import java.util.List;

public abstract class Updateable {
	protected List<Updateable> listeners;
	private Method doUpdateMethod;

	public void update(Object... args) {
    	try {
    		Class[] argTypes = new Class[args.length];
    		for (int i = 0; i < args.length; i++) {
    			argTypes[i] = args[i].getClass();
    		}
	    	doUpdateMethod = Updateable.class.getDeclaredMethod("doUpdate", argTypes);

			synchronized (this) {
				doUpdateMethod.invoke(this, args);
				for (Updateable listener : listeners) {
					listener.update();
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
}
