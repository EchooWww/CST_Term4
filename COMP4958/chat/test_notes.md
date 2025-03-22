1. Under the project root, start the chat server by

```bash
iex --sname server -S mix
```

2. Start the proxy server by

```bash
iex --sname proxy -S mix
```

Ping the server from the proxy server:

```elixir
Node.ping(:"server@Bug-Free")
```

Start the proxy server:

```elixir
Chat.ProxyServer.start_link(6666)
```

You can start multiple proxy servers and ping them to the server at the same time (use different ports).

3. Start the client by

```bash
cd lib/chat
iex client.ex
```

Then start the client by

```elixir
# default to localhost:6666
Chat.Client.start
# OR specify the host and port to the corresponding proxy server
Chat.Client.start("localhost", 7777)
```

Now you can enter the prompts, set and modify the username, and send messages to other clients.
