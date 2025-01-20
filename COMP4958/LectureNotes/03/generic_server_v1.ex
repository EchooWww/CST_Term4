defmodule GenericServer do
  def start(m, state) do
    spawn(fn -> loop(m, state) end)
  end

  def call(pid, msg) do
    send(pid,{self(), msg})
    receive do
      x -> x
    end
  end

  def loop(m, state) do
    receive do
      {from, msg} ->
        {reply, new_state} = m.handle_call(msg, from, state)
        send(from, reply)
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

  def handle_call({:square, x}, _from, state) do
    {x * x, state}
  end
end

defmodule CounterServer do
  def start(n\\0) do
    GenericServer.start(__MODULE__, n)
  end

  def inc(pid) do
    GenericServer.call(pid, :inc)
  end

  def handle_call(:inc, _from, state) do
    {state, state + 1}
  end
end
