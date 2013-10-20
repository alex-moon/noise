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
          io:format("Received ~p~n", [Message]),
          rpop_loop(Consumer)
    end.

receiver(Consumer) ->
    %% {Client, Channel} = Consumer,
    receive
        _ -> 
            %% io:format("Notified! (~p)~n", [Channel]),
            rpop_loop(Consumer),
            receiver(Consumer)
    end.

notifier(Consumer) ->
    {_, Channel} = Consumer,
    {Channel, spawn(fun() ->
        receiver(Consumer)
    end)}.

consume(Consumer) ->
    Notifier = notifier(Consumer),
    Subscriber = subscriber:subscriber(Notifier),
    subscriber:subscribe(Subscriber).