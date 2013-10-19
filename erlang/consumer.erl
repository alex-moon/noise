-module(consumer).
-export([consumer/1, consume/1]).

consumer(Channel) -> 
    {ok, Client} = eredis:start_link(),
    {Client, Channel}.

notifier(Consumer) ->
    {Client, Channel} = Consumer,
    {Channel, spawn(receive
        _ -> 
            io:format("Notified! (~p)~n", [Channel]),
            {ok, Message} = eredis:q(Client, ["RPOP", Channel]),
            io:format("Received ~p~n", [Message])
    end)}.

consume(Consumer) ->
    Notifier = notifier(Consumer),
    Subscriber = subscriber:subscriber(Notifier),
    subscriber:subscribe(Subscriber).