defmodule Arithmetic.Worker do
  use GenServer
  # Client API
  def start() do
    GenServer.start(__MODULE__, nil)
  end

  def square(pid, x) do
    GenServer.call(pid, {:square, x})
  end

  def sqrt(pid, x) do
    GenServer.call(pid, {:sqrt, x}, timeout: 100000) # change the timeout to 100000
  end

  # server impleentation
  @impl true # part of the implementation, tell the compiler to check if that is correctly implemented as the GenServer behaviour
  def init(arg) do
    {:ok, arg}
  end

  @impl true
  def handle_call({:square, x}, _from, state) do
    {:reply, x * x, state}
  end

  @impl true
  def handle_call({:sqrt, x}, _from, state) do
    Process.sleep(10000) # simulate a long running operation
    reply = if x >=0, do: :math.sqrt(x), else: :error
    {:reply, reply, state}
  end
end
