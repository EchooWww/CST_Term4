## Lecture 6

### Nodes in Erlang

- Nodes are Erlang instances that can communicate with each other
- We can create nodes on the same machine or on different machines . If we create one node `iex --sname node1` and another node `iex --sname node2`, we can make them communicate with each other by using the `Node.connect(:node1@hostname)` function, or by `Node.ping(:node1@hostname)` to check if the node is alive.

If the node does not exist, `ping` will return `:pang`. If the node exists, `ping` will return `:pong`.

- `node` gives the name of the current node. `Node.list` gives a list of all nodes that are connected to the current node.

When we spawn a process on another node, the other node will execute the function, but IO will be done on the node that spawned the process (the "group leader" node).

```elixir
Node.spawn(:"node1@Bug-Free", fn -> IO.puts("Hello from #{node}") end) # spawns a process on node1. so the output will be "Hello from node1@Bug-Free"
```

```elixir
f = fn x -> pid = self(); fn -> x * x end end
```

```elixir
Process.register(pid, :square)
```

When we register a process, we can call it by its name. For example, we can call the process `:square` by using `send :square, 5` to send the message `5` to the process `:square`. We can send a message to a process on another node by using `send {:"node1@Bug-Free", :square}, 5`.

We can also register a process globally by using `:global.register_name(:square, pid)`. Then we can get the pid of the process by using `pid = :global.whereis_name(:square)`, and we can send a message to the process by using `send (pid, 5)`.

### Implement a TCP server

We can think of a socket as a file, as everything is a file in Unix. We can open a socket, read from it, and write to it.

```elixir
{:ok, socket} = :gen_tcp.listen(4040, [:binary, packet: :line, active: false])
```

`:gen_tcp` is in the erlang library.

TCP is a stream, we need to tell the server when the message ends. The `:packet` specifies how many bytes for the message number. `:line` specifies that the message ends with a newline character.

Active socket and passive socket: an
