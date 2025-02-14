defmodule Counter.Worker do
  use GenServer

  def start_link(name) do
    GenServer.start_link(__MODULE__, 0, name: via(name))
  end

  defp via(name) do
    {:via, :global, {__MODULE__, name}}
  end

  def inc(name, amt \\ 1) do
    GenServer.cast(via(name), {:inc, amt})
  end

  def value(name) do
    GenServer.call(via(name), :value)
  end

  @impl true
  def init(arg) do
    {:ok, arg}
  end

  @impl true
  def handle_cast({:inc, amt}, {name, value}) do
    {:noreply, {name, value + amt}}
  end

  @impl true
  def handle_call(:value, _from, {_, value} = state) do
    {:reply, value, state}
  end

  @impl true
  def terminate(_reason, state) do
    :ets.insert(@store, state)
    :ok
  end
end
