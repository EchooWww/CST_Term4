## Lecture 4

Erlang is "let it crash" language. It is designed to be fault-tolerant. Instead of trying to handle every possible error, it is better to let the program crash and restart it. We will do it through the supervisor.

Before we look into how to implement a supervisor, let's take a look at the possible failures.

### Errors

- `raise` takes a string as an argument and raises an exception with that string as the reason.

- `exit` takes an atom as an argument and exits the process with that atom as the reason.

  We can also run `Process.exit(self(), :kill)` to kill the current process.

- `throw` takes a value as an argument and throws that value as an exception.

Runtime errors are usually `raise`, `throw` are usually used to control flow.

We can catch those errors:

```elixir
try do
  raise "error"
  catch t, v -> {t, v}
end
# returns {:error, %RuntimeError{message: "error"}} # a struct

try do
  exit(:die)
  catch t, v -> {t, v}
  # returns {:exit, :die}
end

try do
  throw(:throw)
  catch t, v -> {t, v}
  # returns {:throw, :throw}
end
```

There's a `Process.link` function that links two processes together. Eventually, there will be a link set up between all processes in the system.

A process can chooose trap exit signals. When a process traps exit signals, it will receive a message instead of crashing when a linked process crashes. By default, exit signals are not trapped, so the process will crash when a linked process crashes.

Special exit signals: `:normal`, `:kill`. `:normal` kills nothing, `:kill` kills everything, and is untrappable.

When a process is killed, the signal is converted to `:killed`, and other linked processes will receive `:killed` as the reason.

`:erlang.term_to_binary([1, "hello, :hi, {:ok, 5}])` will return a binary representation of the list.

`:erlang.binary_to_term(<<131, 108, 0, 0, 0, 1, 97, 104, 2, 104, 2, 100, 0, 3, 111, 107, 97, 5>>)` will return the original list.

### Supervisor

We can create a mix project with `mix new project_name --sup`, this will generate a supervisor file for us in the `lib` directory called `project_name/application.ex`.

Then inside the `application.ex` file, we can define the supervisor:

````elixir




### Structs

Structs are a way to define a map with a specific set of keys. They are defined with `defstruct` and can be created with `%StructName{}`.

```elixir
defmodule Name do
  defstruct [:first, :last]

  def new(firstname, lastname) do
    %Name{first: firstname, last: lastname}
  end
end

defmodule Student do
  defstruct [id:"", name: %Name{}, scores:%{}]

  def new(id, firstname, lastname, scores \\ %{}) do
    %Student{id: id, name: Name.new(firstname, lastname), scores: scores}
  end
  def parse(data) do
    case String.split(data) do
      [id, first, last | scores] ->
        scores = for [c,s] <- Enum.chunk_every(scores,2), into: %{} do
          {c, String.to_integer(s)}
        end
        {:ok, Student.new(id, first, last, scores)}
      _ -> :error
    end
  end

  def read_file(path) do
    {:ok, content} = File.read(path)
    Enum.map(String.split(content, "\n"), &parse/1)
    |> Enum.filter(&(&1 != :error))
    |> Enum.map(fn {_, x} -> x end)
  end
end
````
