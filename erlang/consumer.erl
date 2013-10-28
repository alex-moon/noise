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

notify_receiver(Consumer) ->
    receive
        Val ->
            io:format("Notified: ~p~n", [Val]),
            rpop_loop(Consumer),
            notify_receiver(Consumer)
    end.

notifier(Consumer) ->
    {_, Channel} = Consumer,
    NotifyReceiver = spawn_link(fun() ->
        notify_receiver(Consumer)
    end),
    {Channel, NotifyReceiver}.

consume(Consumer) ->
    Notifier = notifier(Consumer),
    Subscriber = subscriber:subscriber(Notifier),
    subscriber:subscribe(Subscriber).