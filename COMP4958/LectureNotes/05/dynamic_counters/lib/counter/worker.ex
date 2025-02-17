defmodule Counter.Worker do
  use GenServer

  def start_link(name) do
    GenServer.start_link(__MODULE__, 0, name: via(name))
  end

  def inc(name, amt \\ 1) do
    GenServer.cast(via(name), {:inc, amt})
  end

  def value(name) do
    GenServer.call(via(name), :value)
  end

  defp via(name) do
    {:via, Registry, {Counter.Registry, {__MODULE__, name}}}
  end

  @impl true
  def init(arg) do
    {:ok, arg}
  end

  @impl true
  def handle_cast({:inc, amt}, state) do
    {:noreply, state + amt}
  end

  @impl true
  def handle_call(:value, _from, state) do
    {:reply, state, state}
  end
end
