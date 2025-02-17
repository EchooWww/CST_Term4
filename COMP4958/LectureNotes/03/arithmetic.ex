defmodule Arithmetic.Worker do
  use GenServer

  # client
  def start() do
    GenServer.start(__MODULE__, nil)
  end

  def square(pid, x) do
    GenServer.call(pid, {:square, x})
  end

  def sqrt(pid, x) do
    GenServer.call(pid, {:sqrt, x}, 100000)
  end

  # server (implementation)
  @impl true
  def init(arg) do
    {:ok, arg}
  end

  @impl true
  def handle_call({:square, x}, _from, state) do
    {:reply, x * x, state}
  end

  @impl true
  def handle_call({:sqrt, x}, _from, state) do
    Process.sleep(10000)
    reply = if x >= 0, do: :math.sqrt(x), else: :error
    {:reply, reply, state}
  end
end

