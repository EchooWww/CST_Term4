defmodule Chat.Client do
  def start(host \\ "localhost", port \\ 6666) do
    # Connect to the server
    {:ok, socket} = :gen_tcp.connect(String.to_charlist(host), port, [:binary, active: false, packet: :line])

    # Spawn a process to handle incoming messages
    receiver_pid = spawn(fn -> receiver_loop(socket) end)
    :gen_tcp.controlling_process(socket, receiver_pid)

    # Start the user input loop
    result = input_loop(socket)

    # Ensure socket is closed when done
    :gen_tcp.close(socket)

    result
  end

  defp input_loop(socket) do
    case IO.gets("") do
      :eof ->
        # User pressed Ctrl+D, close connection properly
        IO.puts("Connection closed")
        :ok

      {:error, reason} ->
        IO.puts("Error reading input: #{inspect(reason)}")
        {:error, reason}

      line ->
        # Send user input to server
        :gen_tcp.send(socket, line)
        input_loop(socket)
    end
  end

  defp receiver_loop(socket) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, data} ->
        # Print received message
        IO.write(data)
        receiver_loop(socket)

      {:error, :closed} ->
        IO.puts("Connection closed by server")

      {:error, reason} ->
        IO.puts("Error receiving data: #{inspect(reason)}")
    end
  end
end
