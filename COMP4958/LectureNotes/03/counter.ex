defmodule Counter.Server do
  use GenServer

  # client
  def start(n \\ 0) do
    GenServer.start(__MODULE__, n)
  end

  def inc(pid) do
    GenServer.cast(pid, :inc)
  end

  def value(pid) do
    GenServer.call(pid, :value)
  end

  # server
  @impl true
  def init(arg) do
    {:ok, arg}
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

