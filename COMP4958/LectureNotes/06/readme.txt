To run the counters examples, e.g., the dynamic counters
Assume machine name is elixir`
In global_counters/dynamic, start server:
$ iex --sname bart -S mix
iex(bart@elixir)>

In counter_client/dynamic:
- compile .ex file: elixirc *.ex
- start iex
$ iex --sname homer
iex(homer@elixir)> Node.ping(:bart@elixir)
iex(homer@elixir)> Counter.WorkerSupervisor.start_worker("c1")
iex(homer@elixir)> Counter.Worker.inc("c1")
...

