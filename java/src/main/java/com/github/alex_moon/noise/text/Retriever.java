package com.github.alex_moon.noise.text;

import static java.nio.file.StandardWatchEventKinds.ENTRY_CREATE;
import static java.nio.file.StandardWatchEventKinds.ENTRY_MODIFY;

import java.io.IOException;
import java.nio.file.FileSystems;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.WatchEvent;
import java.nio.file.WatchKey;
import java.nio.file.WatchService;

import com.github.alex_moon.noise.core.Core;

public class Retriever extends Thread {
    private WatchService watcher;
    private WatchKey key;
    private Path textsPath;

    public void run() {
        // watch dir
        try {
            watcher = FileSystems.getDefault().newWatchService();
            textsPath = FileSystems.getDefault().getPath(Core.textsDir);
            key = textsPath.register(watcher, ENTRY_CREATE, ENTRY_MODIFY);
        } catch (IOException e) {
            e.printStackTrace();
        }

        // io loop
        for (;;) {
            try {
                // listen
                key = watcher.take();
            } catch (InterruptedException x) {
                return;
            }

            for (WatchEvent<?> event : key.pollEvents()) {
                // we have new files
                WatchEvent<Path> pathEvent = (WatchEvent<Path>) event;
                Path filename = pathEvent.context(); // context() is the filename
                if (filename != null) {
                    Path file = textsPath.resolve(filename);
                    Text text = createTextFromFile(file);
                    text.update(null);
                }
            }

            key.reset();
        }
    }

    private Text createTextFromFile(Path file) {
        String contents = "";
        try {
            byte[] bytes = Files.readAllBytes(file);
            contents = new String(bytes);
        } catch (IOException e) {
            e.printStackTrace();
        }

        String uuidString = file.toString(); // the filename is the UUID
        return new Text(contents, uuidString);
    }
}