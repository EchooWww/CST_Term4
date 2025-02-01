defmodule Card.Worker do
  use GenServer

  def start_link(_n\\nil) do
    IO.puts("Starting Card.Worker..")
    GenServer.start_link(__MODULE__, nil, name: __MODULE__)
  end

  def new() do
    GenServer.cast(__MODULE__, :new)
  end

  def shuffle() do
    GenServer.cast(__MODULE__, :shuffle)
  end

  def count() do
    GenServer.call(__MODULE__, :count)
  end

  def deal(n\\1) do
    GenServer.call(__MODULE__, {:deal, n})
  end

  @impl true
  def handle_cast(:new, _state) do
    deck = for i <- [2, 3, 4, 5, 6, 7, 8, 9, 10, "J", "Q", "K", "A"],
    j <- [0x2663, 0x2666, 0x2665, 0x2660],
    do: to_string(i) <> to_string([j])
    {:noreply, deck}
  end

  @impl true
  def handle_cast(:shuffle, state) do
    {:noreply, Enum.shuffle(state)}
  end

  @impl true
  def handle_call(:count, _from, state) do
    {:reply, length(state), state}
  end

  @impl true
  def handle_call({:deal, n}, _from, state) when is_integer(n) do
    if n <= 0 or n > length(state) do
      response =
        if n <= 0 do {:error, "Invalid number of cards to deal"}
        else {:error, "Not enough cards in deck"}
        end
      {:reply, response, state}
    else
      {taken, remaining} = Enum.split(state, n)
      response = {:ok, taken}
      {:reply, response, remaining}
    end
  end

  @impl true
  def handle_call({:deal, _n}, _from, _state) do
    raise "The number of cards to deal must be an integer"
  end

  @impl true
  def init(_arg) do
    deck = Card.Store.get() || for i <- [2, 3, 4, 5, 6, 7, 8, 9, 10, "J", "Q", "K", "A"], j <- [0x2663, 0x2666, 0x2665, 0x2660], do: to_string(i) <> to_string([j])
    {:ok, deck}
  end

  @impl true
  def terminate(_reason, state) do
    Card.Store.put(state)
    :ok
  end
end
