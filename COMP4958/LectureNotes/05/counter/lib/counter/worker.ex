defmodule Counter.Worker do
  use GenServer

  def start_link(name) do
    GenServer.start_link(__MODULE__, name, name: name)
  end

  def inc(name, amt \\ 1)  do
    GenServer.cast(name, {:inc, amt})
  end

  def dec(name, amt \\ 1) do
    GenServer.cast(name, {:dec, amt})
  end

  def value(name) do
    GenServer.call(name, :value)
  end

  @impl true
  def init(name) do
    {:ok, {name, Counter.Store.get(name)}}
  end

  @impl true
  def handle_cast({:inc, amt}, {name, state}) do
    {:noreply, {name, state + amt}}
  end

  @impl true
  def handle_cast({:dec, amt}, {name, state}) do
    {:noreply, {name, state - amt}}
  end

  @impl true
  def handle_call(:value, _from, {_, value} = state) do
    {:reply, value, state}
  end

#  @impl true
#  def handle_info(:reset, _state) do
#    {:noreply, 0}
#  end

  @impl true
  def terminate(_reason, {name, state}) do
    Counter.Store.put(name, state)
  end
end


