defmodule Chat.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    port = String.to_integer(System.get_env("PORT") || "6666")

    children = [
      Chat.Server,
      %{
        id: Chat.ProxyServer,
        start: {Chat.ProxyServer, :start_link, [port]},
        type: :worker,
        restart: :permanent
      }
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    :ets.new(Chat.Store, [:named_table, :public])

    opts = [strategy: :one_for_one, name: Chat.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
