defmodule CardServer do
  def start() do
    Process.register(spawn(fn -> loop(new_deck()) end), __MODULE__)
  end

  def new() do
    send(__MODULE__, :new)
  end

  def shuffle() do
    send(__MODULE__, :shuffle)
  end

  def count() do
    send(__MODULE__, {:count, self()})
    receive do x -> x end
  end

  def deal(n\\1) do
    send(__MODULE__, {:deal, n, self()})
    receive do
      x->x
    end
  end

  def loop(deck) do
    receive do
      :new -> loop(new_deck())
      :shuffle -> loop(Enum.shuffle(deck))
      {:count, from} -> send(from, length(deck)); loop(deck)
      {:deal, n, from} ->
        if n <= 0 or n > length(deck) do
          response =
            if n <= 0 do {:error, "Invalid number of cards to deal"}
            else {:error, "Not enough cards in deck"}
            end
          send(from, response)
          loop(deck)
        else
          {taken, remaining} = Enum.split(deck, n)
          send(from, {:ok, taken})
          loop(remaining)
        end
    end
  end

  defp new_deck() do
    for i <- [2, 3, 4, 5, 6, 7, 8, 9, 10, "J", "Q", "K", "A"], j <- [0x2663, 0x2666, 0x2665, 0x2660], do: to_string(i) <> to_string([j])
  end
end
