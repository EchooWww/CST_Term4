defmodule Card.Worker do
  use GenServer
  @store Card.Store

  def start_link(name) do
    IO.puts("Starting Card.Worker..")
    GenServer.start_link(__MODULE__, nil, name: via(name))
  end

  defp new_deck() do
    for i <- [2, 3, 4, 5, 6, 7, 8, 9, 10, "J", "Q", "K", "A"], j <- [0x2663, 0x2666, 0x2665, 0x2660], do: to_string(i) <> to_string([j])
  end

  defp via(name) do
    {:via, Registry, {Card.Registry, {__MODULE__, name}}}
  end

  def new(name) do
    GenServer.cast(via(name), :new)
  end

  def shuffle(name) do
    GenServer.cast(via(name), :shuffle)
  end

  def count(name) do
    GenServer.call(via(name), :count)
  end

  def deal(name, n\\1) do
    GenServer.call(via(name), {:deal, n})
  end

  @impl true
  def handle_cast(:new, _state) do
    deck = new_deck()
    {:noreply, deck}
  end

  @impl true
  def handle_cast(:shuffle, {name, deck}) do
    deck = Enum.shuffle(deck)
    {:noreply, {name, deck}}
  end

  @impl true
  def handle_call(:count, _from, {_, deck} = state) do
    {:reply, length(deck), state}
  end

  @impl true
  def handle_call({:deal, n}, _from, {name, deck} = state) when is_integer(n) do
    if n <= 0 or n > length(deck) do
      response =
        if n <= 0 do {:error, "Invalid number of cards to deal"}
        else {:error, "Not enough cards in deck"}
        end
      {:reply, response, state}
    else
      {taken, remaining} = Enum.split(deck, n)
      response = {:ok, taken}
      {:reply, response, {name, remaining}}
    end
  end

  @impl true
  def handle_call({:deal, _n}, _from, _state) do
    raise "The number of cards to deal must be an integer"
  end

  @impl true
  def init(name) do
    deck =
      case :ets.lookup(@store, name) do
        [{^name, v}] -> v
        [] -> new_deck()
      end
      {:ok, {name, deck}}
    end

  @impl true
  def terminate(_reason, state) do
    :ets.insert(@store, state)
    :ok
  end
end
