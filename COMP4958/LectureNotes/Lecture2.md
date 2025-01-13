# Lecture 2: More about Elixir

## More Data Structures

### Ranges

- A range is a sequence of numbers with a start and an end
- Ranges are created using the `..` operator
- We can convert a range to a list with Enum.to_list/1 (Ranges are also enumerable)

We can define function like

```elixir
defmodule M do
  def f(n\\1), do: n

M.f(2) # 2
```

```elixir
defmodule M fo
  def f(n\\1), do: n
  def f(0), do:-1
  def f(n), do: n
  end
end

M.f(2) # 2
```

### Maps

Maps are defined using the `%{}` and `key=>value`syntax. They are key-value pairs where the key can be any value. Maps are unordered key-value pairs.

```elixir
m = %{:a => 1, 2 => :b, "c" => 3}
m[:a] # 1
m[2] # :b
m["c"] # 3
m["d"] # nil (key not found) nil basically means false
```

Special map: when all the keys in the map are **atoms**, we can use the `key: value` syntax

```elixir
m = %{a: 1, b: 2} # Although a is not written as :a, it is treated as an atom
```

We can also use `Map.get` to get the value of a key in a map.

`Map.put` can be used to add a key-value pair to a map, but instead of modifying the map, it returns a new map with the key-value pair added. (All the data structures in Elixir are **immutable**!)

Update syntax: `%{m | key: value}`, this will update the value of the key in the map (change and assign)

#### Pattern matching with maps

The number of elements does not have to match, but the keys must match.

```elixir
%{a: x} = %{a: 1, b: 2}
x # 1
```

### Keyword Lists

Keyword lists are rarely used in the code, but internally, they're usually used in Elixir.

It is a list of tuples where the first element of the tuple is an atom.

```elixir
l = [{:a, 1}, {:b, 2}]
```

For if-else confitions, it's actually keyword lists used internally.

```elixir
if(1<2, do:3, else:4)
```

This is useful when we call a function with multiple arguments.

```elixir
String.split("Hello, World!", ",", trim: true, parts: 2)
```

## Bit Strings

Exlixir is designed for data communication, so it has bit strings.

### Syntax

A bit string is a sequence of bytes. We can use the `<<>>` syntax to define a bit string.

```elixir
<<35::3, 46::5>> # The last 3 bytes of 35 and the last 5 bytes of 46
```

### Charlists

A charlist is a list of numbers, when each number is a unicode code point.

```elixir
[50,51,52] # ~c"234"
```

### Pattern Matching

```elixir
<<x, y, _::bytes>> = <<35::3, 46::4, 65::7, 66::8>>
```

A string is a binary, so if we call `is_binary("Hello")`, it will return true.

```elixir
<<x, y, _::bytes>> = "hello"
x # 104
y # 101
```

The binary and the string is essentially the same thing.

### Common Operations

```elixir
"hello" <> <<0>> # will return the charlist [104, 101, 108, 108, 111, 0]
```

We can call `to_string` to convert a charlist to a string.

```elixir
to_string([104, 101, 108, 108, 111, 0]) # "hello"
```

`?` is used to get the unicode code point of a character.

When we append a zero byte to a string, it will be treated as a charlist.

```elixir
for i <- 2..10, j <- [0x2660, 0x2665, 0x2666, 0x2663], do:"#{i} <> to_string [j]"

for i <- 2..10, j <- [0x2660, 0x2665, 0x2666, 0x2663], do: to_string [?i, j] # This will not work, as ?i will return the unicode of the character i instead of the variable i

```

### Code Points and Graphemes

String.graphemes will return the number of characters in a string.

String.codepoints will return the unicode code points of a string. (Some languages have multiple code points for a single character)

```elixir
iex(13)> String.codepoints("👩‍🚒")
["👩", "\u200D", "🚒"]
iex(14)> String.length("👩‍🚒")
1
iex(15)> String.graphemes("👩‍🚒")
["👩\u200D🚒"]
```

## Processes

```elixir
f = fn x -> Process.sleep(2000); IO.puts x end
1..10 |> Enum.each(f)
```

In order to run heavy tasks in parallel, we can use processes. There's a scheduler that will run the processes in parallel.

We can call `spawn` to create a new process, which takes a function with no arguments.

```elixir
g = fn x -> spawn(fn -> f.(x) end) end
1..10 |> Enum.each(g)
```

The processes are not the same processes as the operating system processes. They are lightweight and are managed by the Erlang VM. We can create A LOT of processes.

```elixir
spawn(fn -> Process.sleep(1000); :ok end)
# This return the PID of the process
```

### Common Operations

`Process.alive?(pid)` will return true if the process is alive.

`self()` will return the PID of the current process.

We can use `send(pid, message)` to send a message to a process. It will be sent to the mailbox of the process.

`receive do` will receive the message from the mailbox.

```elixir
receive do
  x -> x
end
```

We can do something like:

```elixir
send(self(), :error)
send(self(), {:ok, 42})

receive do
  {:ok, x} -> IO.puts "Result: #{x}"
end
```

`flush()` is a helper function available in the **iex** shell. It prints and clears all messages from the current process’s mailbox.

### Server

We can use these to create a server!

```elixir
defmodule ArithmeticServer do
  def start() do
    spawn(&loop/0) # This will create a new process and run the loop function.
  end

  defp loop() do
    receive do
      {:square, x, from} ->
        send(from, {self(), x * x})
      {:sqrt, x, from} ->
        response =
          if x >= 0, do: {:ok, :math.sqrt(x)}, else: :error
        send(from, response)
      _ -> :ok  # throw away other messages
    end
    loop()
  end
end
```

But this will require the client to call `loop` with the type of operation.

```elixir
pid = ArithmeticServer.start()
send(pid, {:square, 5, self()})
```

This is not very user-friendly. We should have a simplifier interface.

### Client API

The client API is what we add to the server to make it easier to use. It is user to forward the messages to the server, so the users don't need to know the details of the server.

Adding these to the module:

```elixir

  def square(pid, x) do
    send(pid, {:square, x, self()}) # forward the message from client call to the server, i.e., the loop function
    receive do
      {^pid, x} -> x # we only want the message from the server
    end
  end

  def sqrt(pid, x) do
    send(pid, {:sqrt, x, self()})
    receive do
      x -> x
    end
  end
```

Now the client can call the server like this:

```elixir
pid = ArithmeticServer.start()
ArithmeticServer.square(pid, 5)
```

Our `loop` function is now a server that can handle multiple clients. It's only talking to messages forwarded by the client API, instead of the client directly.

### Server that can handle states

We can also, of course, create a server that can handle states. The easiest form is to have a counter

```elixir
defmodule CounterServer do
  def start(n \\ 0) do
    spawn(fn -> loop(n) end)
  end

  def inc(pid) do
    send(pid, :inc)
  end

  def dec(pid) do
    send(pid, :dec)
  end

  def value(pid) do
    send(pid, {:value, self()}) # We're sending the client's PID to the server
    receive do
      x -> x # We're expecting a message from the server
    end
  end

  defp loop(n) do
    receive do
      :inc ->
        loop(n + 1)
      :dec ->
        loop(n - 1)
      {:value, from} ->
        send(from, n) # We only send a message to the client that asked for the value
        loop(n) # Keep looping
    end
  end
end
```

Now we can call the server like this:

```elixir
pid = CounterServer.start()
CounterServer.inc(pid)
CounterServer.inc(pid)
CounterServer.dec(pid)
CounterServer.value(pid) # 1
```

### Register the server

Sometimes it's not very user-friendly to use the PID to call the server. We can register the server with a name.

`Process.register(pid, :name)` will register the process with a name.

send(:name, message) will send the message to the process with the name.

We can call `Process.whereis(:name)` to get the PID of the process with the name.

`__MODULE__` is a macro that will return the name of the current module. This way, we can omit passing the pid to the functions!

```elixir
defmodule RegisteredCounterServer do
  def start(n\\0) do
    Process.register(spwan(fn -> loop(n) end),__MODULE__)
  end

  def int() do
    send(__MODULE__, :inc)
  end

  def dec() do
    send(__MODULE__, :dec)
  end

  def value() do
    send(__MODULE__, {:value, self()})
    receive do
      x -> x
    end
  end

  def loop(n) do
    receive do
      :inc -> loop(n+1)
      :dec -> loop(n-1)
      {:value, from} ->
        send(from, n)
        loop(n)
    end
  end
end
```

## Misc

### List Comprehension

```elixir
for i <- 1..100, rem(i, 3) == 1, do: i
```

### IEx

we can use `v(numerical index)` to get the result of a command in iex.

```elixir
v(1) # [{:a, 1}, {:b, 2}]
```
