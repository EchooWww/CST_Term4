defmodule ArithmeticServer do
  def start() do
    spawn(&loop/0)
  end

  # this version checks that reply is from the server
  def square(pid, x) do
    send(pid, {:square, x, self()})
    receive do
      {^pid, x} -> x
    end
  end

  def sqrt(pid, x) do
    send(pid, {:sqrt, x, self()})
    receive do
      x -> x
    end
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
    send(pid, {:value, self()})
    receive do
      x -> x
    end
  end

  defp loop(n) do
    receive do
      :inc ->
        loop(n + 1)
      :dec ->
        loop(n - 1)
      {:value, from} ->
        send(from, n)
        loop(n)
    end
  end
end

defmodule RegisteredCounterServer do
  def start(n \\ 0) do
    Process.register(spawn(fn -> loop(n) end), __MODULE__)
  end

  def inc() do
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

  defp loop(n) do
    receive do
      :inc ->
        loop(n + 1)
      :dec ->
        loop(n - 1)
      {:value, from} ->
        send(from, n)
        loop(n)
    end
  end
end

