defmodule Counter.Store do
  use GenServer
  def start_link(file) do
    GenServer.start_link(__MODULE__, file, name: __MODULE__)
  end

  def put(value) do
    GenServer.cast(__MODULE__, {:put, value})
  end

  def get do
    GenServer.call(__MODULE__, :get)
  end

  @impl true
  def init(file) do
    {:ok, file}
  end

  @impl true
  def handle_cast({:put, key, value}, file) do
    :ok = File.write(read_file(file), key, value)
    {:noreply, file}
  end

  @impl true
  def handle_call({:get, key}, _from, file) do
    value =
      Map.get(read_file(file), key, 0)
    {:reply, value, file}
  end

  defp read_file(file) do
    case File.read(file) do
      {:ok, content} -> :erlang.binary_to_term(content)
      {:error, _} -> nil
    end
  end
end
