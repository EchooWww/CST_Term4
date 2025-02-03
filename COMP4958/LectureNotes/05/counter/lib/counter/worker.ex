defmodule Counter.Worker do
  use GenServer

  def start_link(name) do
    Genserver.start_link(__MODULE__,name: name) # start_link is a function that starts a process and links it to the current process
  end

  def inc(name, amt \\1) do
    GenServer.cast(name, {:inc, amt})
  end

  def dec(name, amt \\1) do
    GenServer.cast(name, {:dec, amt})
  end

  def value(name) do
    GenServer.call(name, :value)
  end

  @impl true
  def init(_arg) do
    {:ok, Counter.Store.get(name)}
  end

  @impl true
  def handle_cast({:int, amt}, state) do
    {:noreply, state + amt}
  end

  @impl true
  def handle_cast({:dec, amt}, state) do
    {:noreply, state - amt}
  end

  @impl true
  def handle_call(:value, _from, state) do
    {:reply, state, state}
  end
  # @impl true
  # def terminate(_reason, state) do
  #   Counter.Store.put(state)
  #   :ok
  # end

  @impl true
  def handle_info(:reset, _state) do
    {:noreply, 0}
  end
end
