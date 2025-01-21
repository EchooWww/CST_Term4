# Lecture 3

## Generic Servers

### Generic Server and its Implementation

So far, we have written a couple of servers, but they are not generic. A generic server can be designed to implement any server functionality. To make a server generic, we pass the module name to the `start` function.

The generic server handles the loop and message passing (from `call` and `cast` to `handle_call` and `handle_cast`). The callback module provides the specific implementation for the server.

```elixir
defmodule GenericServer do
  def start(m, state) do
    state = m.init(state) # Initialize the state, e.g., open a file if the state is a file.
    spawn(fn -> loop(m, state) end)
  end

  def call(pid, msg) do
    send(pid, {:call, self(), msg})
    receive do
      x -> x
    end
  end

  def cast(pid, msg) do
    send(pid, {:cast, msg})
  end

  def loop(m, state) do
    receive do
      {:call, from, msg} ->
        {reply, new_state} = m.handle_call(msg, from, state)
        send(from, reply)
        loop(m, new_state)

      {:cast, msg} ->
        new_state = m.handle_cast(msg, state)
        loop(m, new_state)
    end
  end
end
```

The implementation includes a Client API and a Server implementation:

```elixir
defmodule CounterServer do
  # Client API
  def start(n \ 0) do
    GenericServer.start(__MODULE__, n)
  end

  def inc(pid) do
    GenericServer.cast(pid, :inc)
  end

  def dec(pid) do
    GenericServer.cast(pid, :dec)
  end

  def value(pid) do
    GenericServer.call(pid, :value)
  end

  # Server-side Implementation (Internal Server, not private because it's used by the generic server)
  def init(arg) do
    arg # We can initialize a more complex state here.
  end

  def handle_cast(:inc, state) do
    state + 1
  end

  def handle_cast(:dec, state) do
    state - 1
  end

  def handle_call(:value, _from, state) do
    {state, state} # The underscore indicates that the variable is unused.
  end
end
```

To initialize the state, the custom server module can implement an `init` function. The flow is as follows:

1. User calls the Client API.
2. The Client API calls the generic server (via `call` or `cast`).
3. The generic server's loop calls the server implementation's `handle_call` or `handle_cast` functions.
4. The handle functions return results to the generic server.
5. The generic server updates the state and continues looping to handle future messages.

### GenServer

Erlang provides the `GenServer` module, which operates similarly to the generic server above. To use it, include `use GenServer` at the top of the module and call `GenServer.start`, `GenServer.cast`, and `GenServer.call` in the Client API.

```elixir
defmodule Counter.Server do
  use GenServer

  # Client API
  def start(n \ 0) do
    GenServer.start(__MODULE__, n)
  end

  def inc(pid) do
    GenServer.cast(pid, :inc)
  end

  def value(pid) do
    GenServer.call(pid, :value)
  end

  # Server-side Implementation
  @impl true # This is a directive to the compiler to ensure the function is implemented.
  def init(arg) do
    {:ok, arg} # The GenServer expects the init function to return a tuple with `:ok` and the state.
  end

  @impl true
  def handle_call(:value, _from, state) do
    {:reply, state, state}
  end

  @impl true
  def handle_cast(:inc, state) do
    {:noreply, state + 1}
  end
end
```

- **`handle_call`**: Returns a tuple with `:reply`, the result to return, and the updated state. Note if we're handling some asynchronous operation, we can return `{:noreply, state}` and send the result later.
- **`handle_cast`**: Returns a tuple with `:noreply` and the updated state.

## Some Other Features of GenServer

### More Concurrent Solutions

When calculations are lengthy, subsequent calculations should not wait. To achieve this, a server can spawn multiple workers to handle requests concurrently. By convention, these are called `Arithmetic.Server` and `Arithmetic.Worker`.

### Tolerance to Failures

The `Supervisor` module can manage and restart workers if they fail, ensuring the system remains robust.

## Streams

Streams are lazy, composable enumerables similar to streams in OCaml. They avoid making copies of lists during transformations, improving efficiency.

Example without streams:

```elixir
1..100_000
|> Enum.map(&(&1 * &1))
|> Enum.filter(fn x -> rem(x, 9) == 1 end)
|> Enum.sum
```

This approach creates copies of the list at each step. Using streams eliminates these copies:

```elixir
1..100_000
|> Stream.map(&(&1 * &1))
|> Stream.filter(fn x -> rem(x, 9) == 1 end)
|> Enum.sum
```

Other stream examples:

```elixir
nats = Stream.iterate(1, &(&1 + 1))
Enum.take(nats, 10)

Stream.cycle([1, 2, 3])
|> Enum.take(10) # [1, 2, 3, 1, 2, 3, 1, 2, 3, 1]
```

FizzBuzz with streams:

```elixir
f = Stream.cycle(["", "", "Fizz"])
b = Stream.cycle(["", "", "", "", "Buzz"])
n = Stream.iterate(1, &(&1 + 1))
fb = Stream.zip_with([f, b, n], fn [f, b, n] -> if f == "" && b == "", do: n, else: f <> b end)
fb |> Enum.take(100)
```

Using comprehensions with streams:

```elixir
for {k, v} <- [a: 1, b: 2, c: 3], into: %{}, do: {k, v} # %{a: 1, b: 2, c: 3}
for {k, 1} <- [a: 1, b: 2, c: 3], do: k # [:a]
```

## Misc Topics

### Create a Project with Mix

To create a new project, run `mix new project_name`. This creates a new directory with the project structure.

Then the convention is to create a folder under the `lib` directory with the same name as the project. This folder will contain the project's modules.

For example, if we want to create a card server, we would make the project structure as follows:

```plaintext
card/
  lib/
    card/
      server.ex
  mix.exs
```

Then, we can run the project under the project root directory with `iex -S mix`.

### Run a Node and Connect from Multiple Shells

To run a process and interact with it from multiple shells, you can use **nodes** in Elixir. Each shell can act as a node, and nodes can communicate with each other.

First, start a node using `iex --sname [name]`. For example:

```bash
iex --sname foo -S mix
# or
iex --sname foo card.ex
```

Then, connect to this node from another shell:

```bash
iex --sname bart --remsh foo@[hostname]
```

Once connected, the second shell (`bart`) interacts directly with the `foo` node, sharing the same processes and state.
