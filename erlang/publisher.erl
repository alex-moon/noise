-module(publisher).
-export([publisher/1, publish/2]).

publisher(Channel) -> 
    {ok, Client} = eredis:start_link(),
    {Client, Channel}.

publish(Publisher, Value) ->
    {Client, Channel} = Publisher,
    [{ok, _}, {ok, _}] = eredis:qp(Client, [
        ["LPUSH", Channel, Value], 
        ["PUBLISH", Channel, 1]
    ]).
