# motivation for GenServer
defmodule GenericServer do
  def start(m, arg) do
    state = m.init(arg)
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

  defp loop(m, state) do
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

defmodule ArithmeticServer do
  def start() do
    GenericServer.start(__MODULE__, nil)
  end

  def square(pid, x) do
    GenericServer.call(pid, {:square, x})
  end

  def init(arg) do
    arg
  end

  def handle_call({:square, x}, _from, state) do
    {x * x, state}
  end
end

defmodule CounterServer do
  # client
  def start(n \\ 0) do
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

  # server
  def init(arg) do
    arg
  end

  def handle_cast(:inc, state) do
    state + 1
  end

  def handle_cast(:dec, state) do
    state - 1
  end

  def handle_call(:value, _from, state) do
    {state, state}
  end
end
