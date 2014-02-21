package com.github.alex_moon.noise.text;

import java.io.IOException;
import java.nio.file.*;
import java.util.UUID;
import static java.nio.file.StandardWatchEventKinds.ENTRY_CREATE;
import static java.nio.file.StandardWatchEventKinds.ENTRY_MODIFY;

public class Retriever extends Thread {
    private final String textsDir = "/home/moona/work/noise/text";
    private WatchService watcher;
    private WatchKey key;
    private Path textsPath;

    public void run() {
        // watch dir
        try {
            watcher = FileSystems.getDefault().newWatchService();
            textsPath = FileSystems.getDefault().getPath(textsDir);
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

            for (WatchEvent<?> event: key.pollEvents()) {
                // we have new files
                WatchEvent<Path> pathEvent = (WatchEvent<Path>)event;
                Path filename = pathEvent.context();  // context() is the filename
                Path file = textsPath.resolve(filename);
                Text text = createTextFromFile(file);
                text.update(null);

                System.out.println("We have a file - " + file + ":");
                for (String line : text.asWordList()) {
                    System.out.println(line);
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

        String uuidString = file.toString();  // the filename is the UUID
        return new Text(contents, uuidString);
    }
}