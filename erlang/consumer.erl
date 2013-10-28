-module(consumer).
-export([consumer/1, consume/1]).

consumer(Channel) -> 
    {ok, Client} = eredis:start_link(),
    {Client, Channel}.

rpop_loop(Consumer) ->
    {Client, Channel} = Consumer,
    {ok, Message} = eredis:q(Client, ["RPOP", Channel]),
    case Message of
        undefined -> undefined;
        _ ->
          io:format("Popped ~p~n", [Message]),
          rpop_loop(Consumer)
    end.

consumer_receiver(Consumer) ->
    receive
        Val ->
            io:format("Notified: ~p~n", [Val]),
            rpop_loop(Consumer),
            consumer_receiver(Consumer)
    end.

notifier(Consumer) ->
    {_, Channel} = Consumer,
    {Channel, spawn_link(fun() ->
        consumer_receiver(Consumer)
    end)}.

consume(Consumer) ->
    Notifier = notifier(Consumer),
    Subscriber = subscriber:subscriber(Notifier),
    subscriber:subscribe(Subscriber).