defmodule ActiveEchoServer do
  require Logger

  def start(port \\ 1234) do
    {:ok, lsocket} = :gen_tcp.listen(port,
      [:binary, active: true, packet: :line, reuseaddr: true])
    accept(lsocket)
  end

  defp accept(lsocket) do
    {:ok, socket} = :gen_tcp.accept(lsocket)
    spawn(fn -> accept(lsocket) end)
    Logger.info("#{inspect self()} connection established: #{inspect socket}")
    loop(socket)
  end

  defp loop(socket) do
    receive do
      {:tcp, ^socket, data} ->
        :gen_tcp.send(socket, data)
        loop(socket)
      {:tcp_closed, ^socket} ->
        Logger.info("#{inspect self()} connection closed: #{inspect socket}")
        :gen_tcp.close(socket)
    end
  end
end

defmodule PassiveEchoServer do
  require Logger

  def start(port \\ 1234) do
    {:ok, lsocket} = :gen_tcp.listen(port,
      [:binary, active: false, packet: :line, reuseaddr: true])
    accept(lsocket)
  end

  defp accept(lsocket) do
    {:ok, socket} = :gen_tcp.accept(lsocket)
    spawn(fn -> accept(lsocket) end)
    Logger.info("#{inspect self()} connection established: #{inspect socket}")
    loop(socket)
  end

  defp loop(socket) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, data} ->
        :gen_tcp.send(socket, data)
        loop(socket)
      {:error, :closed} ->
        Logger.info("#{inspect self()} connection closed: #{inspect socket}")
        :gen_tcp.close(socket)
    end
  end
end

defmodule HybridEchoServer do
  require Logger

  def start(port \\ 1234) do
    {:ok, lsocket} = :gen_tcp.listen(port,
      [:binary, active: false, packet: :line, reuseaddr: true])
    accept(lsocket)
  end

  defp accept(lsocket) do
    {:ok, socket} = :gen_tcp.accept(lsocket)
    spawn(fn -> accept(lsocket) end)
    Logger.info("#{inspect self()} connection established: #{inspect socket}")
    loop(socket)
  end

  defp loop(socket) do
    :inet.setopts(socket, active: :once)
    receive do
      {:tcp, ^socket, data} ->
        :gen_tcp.send(socket, data)
        loop(socket)
      {:tcp_closed, ^socket} ->
        Logger.info("#{inspect self()} connection closed: #{inspect socket}")
        :gen_tcp.close(socket)
    end
  end
end


