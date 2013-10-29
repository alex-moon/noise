-module(subscriber).
-export([subscriber/1, subscribe/2]).

subscriber(Notifier) -> 
    {ok, Client} = eredis_sub:start_link(),
    {Client, Notifier}.

subscribe(Subscriber, Supervisor) ->
    io:format("Subscribing: ~p~n", [Subscriber]),
    {Client, {Channel, _}} = Subscriber,
    process_flag(trap_exit, true),
    spawn_link(fun() ->
        eredis_sub:controlling_process(Client),
        eredis_sub:subscribe(Client, [Channel]),
        subscriber_receiver(Subscriber)
    end),
    receive
        Val ->
            Supervisor ! Val
    end.

subscriber_receiver(Subscriber) ->
    {Client, Notifier} = Subscriber,
    {Channel, NotifyReceiver} = Notifier,
    receive
        {message, BinaryChannel, Message, _} ->
            BinaryChannel = list_to_binary(Channel),
            case Message of
                <<"FuckErlang">> ->
                    io:format("Erlang is fucked - dying...~n");
                _ ->
                    NotifyReceiver ! Message,
                    eredis_sub:ack_message(Client),
                    subscriber_receiver(Subscriber)
            end;
        {subscribed, _, _} ->
            eredis_sub:ack_message(Client),
            subscriber_receiver(Subscriber);
        {eredis_disconnected, _} ->
            exit(Client, kill);
        Val ->
            io:format("PANIC! Expected message or subscribed but got ~p instead~n", [Val]),
            eredis_sub:ack_message(Client),
            subscriber_receiver(Subscriber)
    end.