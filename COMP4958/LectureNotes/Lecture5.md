## Lecture 5: Dynamic Supervisor in Elixir

Starting from the counter example in the previous lecture, we will now implement a dynamic supervisor to manage the counter processes.

When we start the worker, instead of the initial value of the counter, we will pass the name of the counter. This way, we can have multiple registered counters with different names.

```elixir

defmodule Counter.Worker do
  use GenServer

  def start_link(name) do
    Genserver.start_link(__MODULE__, n, name: name) # start_link is a function that starts a process and links it to the current process
  end

  def inc(name, amt \\1) do
    GenServer.cast(name, {:inc, amt})
  end

  def dec(name, amt \\1) do
    GenServer.cast(name, {:dec, amt})
  end

  def value(name) do
    GenServer.call(name, :value)
  end
```

We need to modify the init function and the way data is stored in the state, but we'll deal with this later

Then we can go to `application.ex` and add the workers:

```elixir
  def start(_type, _args) do
    children = [
      # Starts a worker by calling: Counter.Worker.start_link(arg)
      # {Counter.Worker, arg}
      # {Counter.Store, "counter.db"},
      {Counter.Worker, :W1},
      {Counter.Worker, :W2}
    ]
```

But we still get this error: although the workers have different names, the id is not unique. We need to call `Supervisor.child_spec`

```elixir
def start(_type, _args) do
    children = [
      # Starts a worker by calling: Counter.Worker.start_link(arg)
      # {Counter.Worker, arg}
      # {Counter.Store, "counter.db"},
      Supervisor.child_spec({Counter.Worker, :W1}, id: :W1),
      Supervisor.child_spec({Counter.Worker, :W2}, id: :W2)
    ]

```

### Registry

We can use the Registry to register the processes and access them by name. We can use the `Registry` module to register the processes and access them by name.

```elixir
Registry.start_link(name:R, keys: :unique)
```

Then we can register the processes:

```elixir
Registry.register(R, :hello, 1)
```

Then we can look up the process by name:

```elixir
Registry.lookup(R, :hello)
# this will return the PID of the process where the name is registered, so it would be the shell's PID
```

When we supervise processes, we need to use names instead of PIDs as PID can change when the process restarts. We can use the `Registry` module to register the processes and access them by name. More specifically, the `:via` tuple.

```elixir
# 2nd element is the module name
# 3rd element is the tuple with the name of the registry and the key of the record
name = {:via, Registry, {R, :hello}}

# Then when creating a worker, we can pass the name as the first argument
Supervisor.child_spec({Counter.Worker, name}, id: :W1)
```

Then it's time to create a dynamic_counters directory:

```elixir
mix new dymanic_counters --sup --module Counter
```

We can implement the `Counter.Worker` module:

```elixir
defmodule Counter.Worker do
  use GenServer

  def start_link(name) do
    GenServer.start_link(__MODULE__, 0, name: via(name))
  end

  def inc(name, amt\\1) do
    GenServer.cast(via(name), {:inc, amt})
  end

  def value(name) do
    GenServer.call(via(name), :value)
  end

  defp via(name) do
    {:via, Registry, {Counter.Registry, {__MODULE__, name}}}
  end

  @impl true
  def init(arg) do
    {:ok, arg}
  end

  @impl true
  def handle_cast({:inc, amt}, state) do
    {:noreply, state + amt}
  end

  @impl true
  def handle_call(:value, _from, state) do
    {:reply, state, state}
  end
end

```

Note we have a helper function `via` that returns the tuple with the registry name and the key. We also have the `Counter.WorkerSupervisor` module that we need to implement:

```elixir
defmodule Counter.WorkerSupervisor do
  use DynamicSupervisor

  def start_link(_) do
    DynamicSupervisor.start_link(__MODULE__, nil, name: __MODULE__)
  end


  def start_worker(name) do
    DynamicSupervisor.start_child(__MODULE__, {Counter.Worker, name}) # this is equivalent to Supervisor.child_spec({Counter.Worker, name}), which calls Counter.Worker.start_link(name)
  end

  def init(_arg) do
    DynamicSupervisor.init(strategy: :one_for_one, max_children: 100) # if one worker dies, start one worker
  end

end
```

In this module, we need to give the supervisor a name and implement the `init` function. We also have the `start_worker` function that starts a worker process.

Then in `application.ex` we need to add the registry, and the supervisor to the children list:

```elixir
  def start(_type, _args) do
    children = [
      {Registry, name: Counter.Registry, keys: :unique},
      Counter.WorkerSupervisor
    ]
```

Then we can start the supervisor when running the application:

```elixir
iex -S mix
```

Then we can start a worker:

```elixir
Counter.WorkerSupervisor.start_worker(:W1)
```

Now we wanna store the state when the process dies, we're gonna use ets: Erlang Term Storage. We can use the `:ets` module to create a table and store the state in it. Erlang also has `dets` and `mnesia` for disk storage.

```elixir
r = :ets.new(T, []) # returns a reference
:ets.insert(r, {1, "hello"}) # returns true for storing a key-value pair
:ets.lookup(r, 1) # returns [{1, "hello"}]
```

We can also give the table a name instead of a reference:

```elixir
r = :ets.new(Tab, [:named_table, :public]) # now it returns the name of the table
:ets.insert(Tab, {1, "hello"}) # a key can be anything, even a list
```

We have a powerful `match_object` function that can match a pattern:

```elixir
:ets.match_object(Tab,{:_, 1}) # returns all the pairs with the value 1
```

We can also add a `:bag` option to allow duplicate keys:

```elixir
r = :ets.new(Tab, [:named_table, :public, :bag])
:ets.insert(Tab, {1, "hello"})
:ets.insert(Tab, {1, "world"})
:ets.lookup(Tab, 1) # returns [{1, "hello"}, {1, "world"}]
```

but if the whole tuple is the same, it will not be stored:

```elixir
:ets.insert(Tab, {1, "hello"})
:ets.insert(Tab, {1, "hello"})
:ets.lookup(Tab, 1) # returns [{1, "hello"}]
```

So we can use `:ets` to store the state of the counter process. We can look up the state in the `init` function with ets:

```elixir
  @impl true
  def init(name) do
    value =
    case :ets.lookup(@store, name) do
      [{^name, v}] -> v
      [] -> 0
    end
    {:ok, {name, value}}
  end
```

and store them in `terminate`:

```elixir
  @impl true
  def terminate(_reason, state) do
    :ets.insert(@store, state)
    :ok
  end
```

Then it's fairly easy to start ets in `application.ex`:

```elixir
  def start(_type, _args) do
    children = [
      {Registry, name: Counter.Registry, keys: :unique},
      Counter.WorkerSupervisor,
      {Counter.Store, name: Counter.Store}
    ]
    :ets.new(Counter.Store, [:named_table, :public])
  end
```

Ets is quite robust, so it doesn't need to be supervised most of the time.
