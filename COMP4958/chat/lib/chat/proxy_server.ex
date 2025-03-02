defmodule Chat.ProxyServer do
  require Logger
  use Task

  @default_port 6666
  def start_link(port \\ @default_port) do
    Task.start_link(fn -> accept(port) end)
  end

  def accept(port) do
    # Open the socket in passive mode, handling potential errors
    case :gen_tcp.listen(port, [
      :binary,
      packet: :line,
      active: false,
      reuseaddr: true
    ]) do
      {:ok, socket} ->
        Logger.info("ProxyServer accepting connections on port #{port}")
        loop_acceptor(socket)

      {:error, reason} ->
        Logger.error("Failed to start proxy server: #{inspect(reason)}")
    end
  end

  defp loop_acceptor(socket) do
    # Accept client connections
    {:ok, client} = :gen_tcp.accept(socket)

    # Spawn a process to handle this client
    {:ok, pid} = Task.start_link(fn -> serve(client) end)
    :ok = :gen_tcp.controlling_process(client, pid)

    # Continue accepting connections
    loop_acceptor(socket)
  end

  defp serve(socket) do
    # Initial message to client
    :gen_tcp.send(socket, "Welcome to the Elixir Chat Server!\r\n")
    :gen_tcp.send(socket, "Please set a nickname with /NICK <nickname> or /N <nickname>\r\n")

    # Configure socket for active mode - better for performance
    :inet.setopts(socket, [active: true])

    # Start processing client messages
    socket_loop(socket, nil)
  end

  defp socket_loop(socket, nickname) do
    receive do
      # Handle incoming chat messages from other users
      {:chat_message, sender, message} ->
        formatted_message = "#{sender}: #{message}\r\n"
        :gen_tcp.send(socket, formatted_message)
        socket_loop(socket, nickname)

      # Handle TCP data
      {:tcp, ^socket, data} ->
        # Process the command
        command = String.trim(data)
        {new_nickname, response} = process_command(command, nickname)

        # Send response back to client
        :gen_tcp.send(socket, response <> "\r\n")

        # Continue the loop with potentially updated nickname
        socket_loop(socket, new_nickname)

      # Handle TCP closure
      {:tcp_closed, ^socket} ->
        # Client disconnected, clean up nickname if registered
        if nickname, do: Chat.Server.unregister_nickname(nickname)
        Logger.info("Client disconnected#{if nickname, do: " (#{nickname})", else: ""}")
        :ok

      # Handle TCP errors
      {:tcp_error, ^socket, reason} ->
        Logger.error("Socket error: #{inspect(reason)}")
        if nickname, do: Chat.Server.unregister_nickname(nickname)
        :ok
    end
  end

  defp process_command(command, nickname) do
    cond do
      # NICK command - exactly as specified in the requirements
      String.match?(command, ~r{^/NICK\s+\S+}) or String.match?(command, ~r{^/N\s+\S+}) ->
        new_nickname = command
                      |> String.split(" ", parts: 2)
                      |> List.last
                      |> String.split(" ")
                      |> List.first

        case Chat.Server.register_nickname(new_nickname, self()) do
          {:ok, message} ->
            # If changing nickname, unregister old one
            if nickname && nickname != new_nickname do
              Chat.Server.unregister_nickname(nickname)
            end
            {new_nickname, message}

          {:error, message} ->
            {nickname, message}
        end

      # LIST command - no arguments needed
      String.match?(command, ~r{^/LIST$}) or String.match?(command, ~r{^/L$}) ->
        case Chat.Server.list_users() do
          {:ok, []} ->
            {nickname, "No users currently connected."}

          {:ok, users} ->
            users_str = Enum.join(users, ", ")
            {nickname, "Users: #{users_str}"}
        end

      # Allow LIST with ignored arguments
      String.match?(command, ~r{^/LIST\s}) or String.match?(command, ~r{^/L\s}) ->
        case Chat.Server.list_users() do
          {:ok, []} ->
            {nickname, "No users currently connected."}

          {:ok, users} ->
            users_str = Enum.join(users, ", ")
            {nickname, "Users: #{users_str}"}
        end

      # MSG command - requires recipient and message
      String.match?(command, ~r{^/MSG\s+\S+\s+\S+}) or String.match?(command, ~r{^/M\s+\S+\s+\S+}) ->
        if nickname do
          # Parse recipient and message
          [_cmd | rest] = String.split(command, " ", parts: 2)
          [recipients_str | message_parts] = String.split(List.to_string(rest), " ", parts: 2)
          message = List.first(message_parts) || ""

          # Send the message
          Chat.Server.send_message(nickname, recipients_str, message)
          recipients = if recipients_str == "*" do
            "all users"
          else
            recipients_str
          end
          {nickname, "Message sent to #{recipients}"}
        else
          {nickname, "You must set a nickname before sending messages. Use /NICK <nickname>"}
        end

      # Invalid MSG command (missing message or recipients)
      String.match?(command, ~r{^/MSG}) or String.match?(command, ~r{^/M}) ->
        {nickname, "Invalid MSG format. Use /MSG <recipients> <message>"}

      # Invalid NICK command (missing nickname)
      String.match?(command, ~r{^/NICK$}) or String.match?(command, ~r{^/N$}) ->
        {nickname, "Invalid NICK format. Use /NICK <nickname>"}

      # Any other command is invalid
      true ->
        {nickname, "Unknown command. Available commands: /NICK <nickname>, /LIST, /MSG <recipients> <message>"}
    end
  end
end
