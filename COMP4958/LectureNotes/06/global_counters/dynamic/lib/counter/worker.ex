defmodule Counter.Worker do
  use GenServer
  @store  Counter.Store

  def start_link(name) do
    GenServer.start_link(__MODULE__, name, name: via(name))
  end

  def inc(name, amt \\ 1) do
    GenServer.cast(via(name), {:inc, amt})
  end

  def value(name) do
    GenServer.call(via(name), :value)
  end

  defp via(name) do
    {:via, :global, {__MODULE__, name}}
  end

  def init(name) do
    value =
      case :ets.lookup(@store, name) do
        [{^name, v}] -> v
        _ -> 0
      end
    {:ok, {name, value}}
  end

  def handle_cast({:inc, amt}, {name, value}) do
    {:noreply, {name, value + amt}}
  end

  def handle_call(:value, _from, {_, value} = state) do
    {:reply, value, state}
  end

  def terminate(_reason, state) do
    :ets.insert(@store, state)
  end
end
