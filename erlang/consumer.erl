-module(consumer).
-export([consumer/1, consume/1]).

consumer(Channel) -> 
    {ok, Client} = eredis:start_link(),
    {Client, Channel}.

receiver(Consumer) ->
    {Client, Channel} = Consumer,
    receive
        _ -> 
            io:format("Notified! (~p)~n", [Channel]),
            {ok, Message} = eredis:q(Client, ["RPOP", Channel]),
            io:format("Received ~p~n", [Message]),
            receiver(Consumer)
    end.

notifier(Consumer) ->
    {_, Channel} = Consumer,
    {Channel, spawn_link(fun() ->
        receiver(Consumer)
    end)}.

consume(Consumer) ->
    Notifier = notifier(Consumer),
    Subscriber = subscriber:subscriber(Notifier),
    subscriber:subscribe(Subscriber).