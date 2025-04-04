defmodule Card.Server do
  use GenServer

  def start() do
    GenServer.start(__MODULE__, nil)
  end

  def new(pid) do
    GenServer.cast(pid, :new)
  end

  def shuffle(pid) do
    GenServer.cast(pid, :shuffle)
  end

  def count(pid) do
    GenServer.call(pid, :count)
  end

  def deal(pid, n\\1) do
    GenServer.call(pid, {:deal, n})
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
  def handle_call({:deal, n}, _from, state) do
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
  def init(_arg) do
    deck = for i <- [2, 3, 4, 5, 6, 7, 8, 9, 10, "J", "Q", "K", "A"], j <- [0x2663, 0x2666, 0x2665, 0x2660], do: to_string(i) <> to_string([j])
    {:ok, deck}
  end
end
