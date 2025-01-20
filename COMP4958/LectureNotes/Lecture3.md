# Lecture 3

## Generic Servers

We have written a couple of servers so far, but they are not generic. We can write a generic server that can be used to implement any server.

To make it generic, we need to pass the module name to the `start` function. This way, we can use the same server to implement different servers.

We also gonna pass a state to the `start` function. This way, we can use the same server to implement different servers with different states.

We can append a underscore to the beginning of a variable name to avoid the warning of unused variables.

Erlang actually provides a `gen_server` module that can be used to implement a generic server. We can use the `gen_server` module to implement a generic server.

## More Concurrent Solutions

If the calculation is lengthy, we don't want the subsequent calculations to wait. Instead, we start a server which starts a number of workers

By convention, we call them `Arithmetic.Server` and `Arithmetic.Worker`.

## Tolerance to Failures

We can use the `supervisor` module to implement a supervisor that can restart the workers if they fail.

## Streams

Lazy composable emurables. Just like streams in OCaml.

```elixir
1..100_000
|> Enum.map(&(&1 * &1))
|> Enum.filter(fn x -> rem(x, 9) == 1 end)
|> Enum.sum
```

In the example above, we're making copies of the list in each step. This is not efficient. We can use streams to avoid making copies of the list.

The idea of using streams is to avoid making copies of the list. We can use the `Stream` module to implement streams. Instead of making copies of the list, it doesn't compute the next element until it's needed.

```elixir
1..100_000
|> Stream.map(&(&1 * &1))
|> Stream.filter(fn x -> rem(x, 9) == 1 end)
|> Enum.sum
```

```elixir
nats = Stream.iterate(1, &(&1 + 1))
Enum.take(nats, 10)
```

```elixir
Stream.cycle([1, 2, 3])
|> Enum.take(10) # [1, 2, 3, 1, 2, 3, 1, 2, 3, 1]
```

We can implement a fizzbuzz with cycles.

```elixir
f = Stream.cycle(["", "", "Fizz"])
b = Stream.cycle(["", "", "", "", "Buzz"])
n = Stream.iterate(1, &(&1 + 1))
fb = Stream.zip_with([f, b, n], fn [f, b, n] -> if f == "" && b == "", do: n, else: f <> b end)
fb |> Enum.take(100)
```

```elixir
for {k, v} <- [a: 1, b: 2, c: 3], into: %{}, do: {k,v} # %{a: 1, b: 2, c: 3}
for {k, 1} <- [a: 1, b: 2, c: 3], do: k # [:a]
```
