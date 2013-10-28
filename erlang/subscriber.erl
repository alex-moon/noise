-module(subscriber).
-export([subscriber/1, subscribe/1]).

subscriber(Notifier) -> 
    {ok, Client} = eredis_sub:start_link(),
    {Client, Notifier}.

subscribe(Subscriber) ->
    io:format("Subscribing: ~p~n", [Subscriber]),
    {Client, {Channel, _}} = Subscriber,
    process_flag(trap_exit, true),
    spawn_link(fun() ->
        eredis_sub:controlling_process(Client),
        eredis_sub:subscribe(Client, [Channel]),
        subscriber_receiver(Subscriber)
    end).

subscriber_receiver(Subscriber) ->
    {Client, {Channel, NotifyReceiver}} = Subscriber,
    receive
        {'EXIT', _, _} ->
            io:format("Our Redis connection has died! Restarting...~n"),
            NotifyReceiver ! "LOL Redis connection died",
            subscribe(Subscriber);
        Val ->
            io:format("Subscriber for ~p received val ~p~n", [Channel, Val]),
            NotifyReceiver ! Val,
            eredis_sub:ack_message(Client),
            subscriber_receiver(Subscriber)
    end.