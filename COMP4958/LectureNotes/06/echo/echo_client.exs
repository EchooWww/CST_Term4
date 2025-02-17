defmodule EchoClient do
  def connect(host, port) do
    {:ok, socket} = :gen_tcp.connect(host, port,
      [:binary, active: false, packet: :line])
    loop(socket)
  end

  defp loop(socket) do
    case IO.gets("> ") do
      :eof ->
        :gen_tcp.close(socket)
      data ->
        :gen_tcp.send(socket, data)
        {:ok, reply} = :gen_tcp.recv(socket, 0)
        IO.write(reply)
        loop(socket)
    end
  end
end

with [host, port | _] <- System.argv(),
     {port, ""} <- Integer.parse(port)
do
  EchoClient.connect(String.to_atom(host), port)
else
  _ -> EchoClient.connect(:localhost, 1234)
end

