package com.github.alex_moon.noise.text;

import java.util.UUID;

public class Text {
    private String stringValue;
    private UUID uuid;

    public Text(String initialStringValue) {
        stringValue = initialStringValue;
        uuid = UUID.randomUUID();
    }

    public Text(String initialStringValue, String uuid) {
        stringValue = initialStringValue;
        try {
            this.uuid = UUID.fromString(uuid);
        } catch (IllegalArgumentException e) {
            this.uuid = UUID.randomUUID();
        }
    }

    public String[] asWordList() {
        String result = stringValue.replaceAll("[^a-zA-Z0-9]+", " ");
        result = result.toLowerCase();
        return result.split(" ");
    }

    public String toString() {
        return stringValue;
    }

    public UUID getUuid() {
        return uuid;
    }

    public String getUuidAsString() {
        return uuid.toString();
    }
}