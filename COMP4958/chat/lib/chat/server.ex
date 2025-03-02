defmodule Chat.Server do
  use GenServer
  require Logger
  @store Chat.Store

  # Client API
  def start_link(_opts) do
    GenServer.start_link(__MODULE__, nil, name: __MODULE__)
  end

  def register_nickname(nickname, pid) do
    GenServer.call(__MODULE__, {:register_nickname, nickname,pid})
  end

  def unregister_nickname(nickname) do
    GenServer.call(__MODULE__, {:unregister_nickname, nickname, self()})
  end

  def list_users do
    GenServer.call(__MODULE__, :list_users)
  end

  def send_message(sender, recipient, message) do
    GenServer.cast(__MODULE__, {:send_message, sender, recipient, message})
  end

  # Server callbacks
  @impl true
  def init(_opts) do
    Logger.info("Starting Chat.Server")
    # Recover nicknames from ETS table
    nicknames = case :ets.lookup(@store, :nicknames) do
      [{:nicknames, nicknames}] -> nicknames
      [] -> %{}
    end
    Enum.each(nicknames, fn {_nickname, pid} ->
      Process.monitor(pid)
    end)
    Logger.info("Chat server initialized with #{map_size(nicknames)} existing users")
    {:ok, %{nicknames: nicknames}}
  end

  @impl true
  def terminate(_reason, state) do
    # Save state to ETS table before shutting down
    :ets.insert(@store, {:nicknames, state.nicknames})
    Logger.info("Chat server terminated, saved #{map_size(state.nicknames)} nicknames to ETS")
    :ok
  end

  @impl true
  def handle_call({:register_nickname, nickname, pid}, _from, state) do
    case validate_nickname(nickname) do
      :ok ->
        if Map.has_key?(state.nicknames, nickname) do
          {:reply, {:error, "Nickname #{nickname} already in use"}, state}
        else
          Process.monitor(pid)
          new_nicknames = Map.put(state.nicknames, nickname, pid)
          Logger.info("User registered with nickname: #{nickname}")
          {:reply, {:ok, "Nickname #{nickname} registered successfully"}, %{state | nicknames: new_nicknames}}
        end
      {:error, reason} ->
        {:reply, {:error, reason}, state}
      end
  end

  @impl true
  def handle_call(:list_users, _from, state) do
    users = Map.keys(state.nicknames) |> Enum.sort()
    {:reply, {:ok, users}, state}
  end

  @impl true
  def handle_call({:unregister_nickname, nickname, _pid}, _from, state) do
    new_nicknames = Map.delete(state.nicknames, nickname)
    Logger.info("Unregistered nickname: #{nickname}")
    {:reply, :ok, %{state | nicknames: new_nicknames}}
  end

  @impl true
  def handle_cast({:send_message, sender, recipients, message}, state) do
    # Handle different recipient types
    recipients_list = parse_recipients(recipients, state.nicknames, sender)
    Enum.each(recipients_list, fn recipient ->
      case Map.get(state.nicknames, recipient) do
        nil ->
          Logger.info("No recipient found with nickname: #{recipient}")

        pid ->
          if Process.alive?(pid) do
            send(pid, {:chat_message, sender, message})
          else
            Logger.warning("Recipient #{recipient} process is dead, will be cleaned up by monitor")
          end
      end
    end)

    {:noreply, state}
  end

  @impl true
  def handle_info({:DOWN, _ref, :process, pid, reason}, state) do
    {nickname, new_nicknames} = find_and_remove_by_pid(state.nicknames, pid)

    if nickname do
      Logger.info("User #{nickname} disconnected: #{inspect(reason)}")
    end

    {:noreply, %{state | nicknames: new_nicknames}}
  end

  defp validate_nickname(nickname) do
    cond do
      # Check if nickname is valid (first char is a letter, rest alphanumeric or underscore, max 12 chars)
      not String.match?(nickname, ~r/^[a-zA-Z][a-zA-Z0-9_]{0,11}$/) ->
        {:error, "Invalid nickname format"}
      true ->
        :ok
    end
  end

  defp parse_recipients("*", nicknames, sender) do
    # Send to all users except the sender
    Map.keys(nicknames) |> Enum.reject(&(&1 == sender))
  end

  defp parse_recipients(recipients, _nicknames, _sender) do
    recipients
    |> String.split(",", trim: true)
    |> Enum.map(&String.trim/1)
  end

  defp find_and_remove_by_pid(nicknames, target_pid) do
    case Enum.find(nicknames, fn {_nick, pid} -> pid == target_pid end) do
      nil ->
        {nil, nicknames}

      {nickname, _} ->
        {nickname, Map.delete(nicknames, nickname)}
    end
  end
end
