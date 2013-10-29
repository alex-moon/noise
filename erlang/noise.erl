-module(noise).
-export([noise/0]).

noise() -> 
    Consumer = consumer:consumer("noise"),
    Publisher = publisher:publisher("noise"),
    spawn_link(fun() -> consumer:consume(Consumer) end),
    [publisher:publish(Publisher, Value) || Value <- ["Oh shit Erlang is speaking to Redis", "Redis is speaking to Go", "you know what this means...", "MAMAAA JUST KILLED A MAN"]].
