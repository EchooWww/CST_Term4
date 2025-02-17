defmodule Counter.Store do
  use GenServer

  def start_link(file) do
    GenServer.start_link(__MODULE__, file, name: __MODULE__)
  end

  def put(key, value) do
    GenServer.cast(__MODULE__, {:put, key, value})
  end

  def get(key) do
    GenServer.call(__MODULE__, {:get, key})
  end

  @impl true
  def init(file) do
    {:ok, file}
  end

  @impl true
  def handle_cast({:put, key, value}, file) do
    m = Map.put(read_file(file), key, value)
    :ok = File.write(file, :erlang.term_to_binary(m))
    {:noreply, file}
  end

  @impl true
  def handle_call({:get, key}, _from, file) do
    value = Map.get(read_file(file), key, 0)
    {:reply, value, file}
  end

  defp read_file(file) do
    case File.read(file) do
      {:ok, content} -> :erlang.binary_to_term(content)
      {:error, _} -> %{}
    end
  end
end

