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
          io:format("~p~n", [binary_to_list(Message)]),
          rpop_loop(Consumer)
    end.

notify_receiver(Consumer) ->
    receive
        _ ->
            rpop_loop(Consumer),
            notify_receiver(Consumer)
    end.

notifier(Consumer) ->
    {_, Channel} = Consumer,
    NotifyReceiver = spawn_link(fun() ->
        notify_receiver(Consumer)
    end),
    {Channel, NotifyReceiver}.

supervise(Subscriber, Consumer) ->
    {Client, _} = Subscriber,
    subscriber:subscribe(Subscriber, self()),
    receive
        {'EXIT', _, _} ->
            case is_process_alive(Client) of
                true ->
                    supervise(Subscriber, Consumer);
                false ->
                    consume(Consumer)
            end;
        Val ->
            io:format("SUPERVISOR: Expected EXIT but got ~p instead~n", [Val])
    end.

consume(Consumer) ->
    Notifier = notifier(Consumer),
    Subscriber = subscriber:subscriber(Notifier),
    process_flag(trap_exit, true),
    supervise(Subscriber, Consumer).