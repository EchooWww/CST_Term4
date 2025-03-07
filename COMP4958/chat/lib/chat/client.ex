defmodule Chat.Client do
  def start(host \\ "localhost", port \\ 6666) do
    # Connect to the server
    {:ok, socket} = :gen_tcp.connect(String.to_charlist(host), port, [:binary, active: false, packet: :line])

    # Spawn a process to handle incoming messages
    main_pid = self()
    receiver_pid = spawn(fn -> receiver_loop(socket, main_pid) end)
    :gen_tcp.controlling_process(socket, receiver_pid)

    # Start the user input loop
    result = input_loop(socket, receiver_pid)

    # Ensure socket is closed when done
    :gen_tcp.close(socket)

    result
  end

  defp input_loop(socket, receiver_pid) do
    case IO.gets("") do
      :eof ->
        # User pressed Ctrl+D, terminate the receiver process first
        Process.exit(receiver_pid, :kill)
        # Close connection
        :gen_tcp.close(socket)
        IO.puts("Connection closed by user")
        :ok
        :closed_by_server ->
          Process.exit(receiver_pid, :kill)

        {:error, reason} ->
          Process.exit(receiver_pid, :normal)
          :gen_tcp.close(socket)
          IO.puts("Error reading input: #{inspect(reason)}")
          {:error, reason}

        line ->
          # Send user input to server
          :gen_tcp.send(socket, line)
          input_loop(socket, receiver_pid)
    end
  end

  defp receiver_loop(socket, main_pid) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, data} ->
        # Print received message
        IO.write(data)
        receiver_loop(socket, main_pid)

      {:error, :closed} ->
        IO.puts("Connection closed by server")
        send(main_pid, :closed_by_server)

      {:error, reason} ->
        IO.puts("Error receiving data: #{inspect(reason)}")
    end
  end
end
