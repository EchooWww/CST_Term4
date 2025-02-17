defmodule Counter.Worker do
  use GenServer
  @name {:global, __MODULE__}

  def start_link(n \\ 0) do
    GenServer.start_link(__MODULE__, n, name: @name)
  end

  def inc(amt \\ 1)  do
    GenServer.cast(@name, {:inc, amt})
  end

  def dec(amt \\ 1) do
    GenServer.cast(@name, {:dec, amt})
  end

  def value() do
    GenServer.call(@name, :value)
  end

  @impl true
  def init(arg) do
    {:ok, 
    case Counter.Store.get() do
      nil -> arg
      value -> value
    end}
  end

  @impl true
  def handle_cast({:inc, amt}, state) do
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

  @impl true
  def handle_info(:reset, _state) do
    {:noreply, 0}
  end

  @impl true
  def terminate(_reason, state) do
    Counter.Store.put(state)
  end
end


