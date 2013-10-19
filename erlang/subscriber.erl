-module(subscriber).
-export([subscriber/1, subscribe/1]).

subscriber(Notifier) -> 
    {ok, Client} = eredis_sub:start_link(),
    {Client, Notifier}.

subscribe(Subscriber) ->
    {Client, {Channel, Notify}} = Subscriber,
    spawn_link(fun () ->
        eredis_sub:controlling_process(Client),
        eredis_sub:subscribe(Client, [Channel]),
        receiver(Notify)
    end).

receiver(Notify) ->
    receive
        Val ->
            Notify ! Val
    end.