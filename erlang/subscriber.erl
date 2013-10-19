-module(subscriber).
-export([subscriber/1, subscribe/1]).

subscriber(Channel) -> 
    {ok, Client} = eredis_sub:start_link(),
    {Client, Channel}.

subscribe(Subscriber) ->
    {Client, Channel} = Subscriber,
    Receiver = spawn_link(fun () ->
        eredis_sub:controlling_process(Client),
        eredis_sub:subscribe(Client, [Channel]),
        receiver(Channel)
    end),
    {Client, Receiver}.

receiver(Channel) ->
    receive
        _ ->
            io:format("Notify consumer ~p~n", [Channel])
    end.