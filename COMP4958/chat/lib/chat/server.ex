defmodule Chat.Server do
  use GenServer
  require Logger
  @store Chat.Store

  # Client API
  def start_link(_opts) do
    GenServer.start_link(__MODULE__, nil, name: {:global, __MODULE__})
  end

  def register_nickname(nickname, pid) do
    GenServer.call({:global, __MODULE__}, {:register_nickname, nickname,pid})
  end

  def unregister_nickname(nickname, reason \\ :disconnect) do
    GenServer.call({:global, __MODULE__}, {:unregister_nickname, nickname, self(), reason})
  end

  def change_nickname(old_nickname, new_nickname, pid) do
    GenServer.call({:global, __MODULE__}, {:change_nickname, old_nickname, new_nickname, pid})
  end

  def list_users do
    GenServer.call({:global, __MODULE__}, :list_users)
  end

  def send_message(sender, recipient, message) do
    GenServer.call({:global, __MODULE__}, {:send_message, sender, recipient, message})
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
  def handle_call({:unregister_nickname, nickname, _pid, reason}, _from, state) do
    new_nicknames = Map.delete(state.nicknames, nickname)

    case reason do
      :disconnect ->
        Logger.info("User #{nickname} left the chat")
      _ ->
        Logger.info("Unregistered nickname: #{nickname}")
    end

    {:reply, :ok, %{state | nicknames: new_nicknames}}
  end

  @impl true
  def handle_call({:change_nickname, old_nickname, new_nickname, pid}, _from, state) do
    case validate_nickname(new_nickname) do
      :ok ->
        if Map.has_key?(state.nicknames, new_nickname) do
          {:reply, {:error, "Nickname #{new_nickname} already in use"}, state}
        else
          new_nicknames = state.nicknames
                          |> Map.delete(old_nickname)
                          |> Map.put(new_nickname, pid)

          Logger.info("User #{old_nickname} changed nickname to #{new_nickname}")
          {:reply, {:ok, "Nickname changed to #{new_nickname} successfully"}, %{state | nicknames: new_nicknames}}
        end
      {:error, reason} ->
        {:reply, {:error, reason}, state}
    end
  end

  @impl true
def handle_call({:send_message, sender, recipients, message}, _from, state) do
  recipients_list = parse_recipients(recipients, state.nicknames, sender)

  {success, failed} = Enum.reduce(recipients_list, {[], []}, fn recipient, {success, failed} ->
    case Map.get(state.nicknames, recipient) do
      nil ->
        Logger.info("No recipient found with nickname: #{recipient}")
        {success, [recipient | failed]}

      pid ->
        if is_pid(pid) and Node.ping(node(pid)) != :pang do
          send(pid, {:chat_message, sender, message})
          {[recipient | success], failed}
        else
          Logger.warning("Recipient #{recipient} process is dead or unreachable, will be cleaned up by monitor")
          {success, [recipient | failed]}
        end
    end
  end)

  {:reply, {Enum.reverse(success), Enum.reverse(failed)}, state}
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
