-module(noise).
-export([noise/0]).



noise() -> 
    {ok, Client} = eredis:start_link(),
    {ok, _} = eredis:q(Client, ["LPUSH", "noise", "Here is a test message"]).