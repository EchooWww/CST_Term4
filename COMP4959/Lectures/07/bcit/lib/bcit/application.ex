defmodule Bcit.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      BcitWeb.Telemetry,
      Bcit.Repo,
      {DNSCluster, query: Application.get_env(:bcit, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: Bcit.PubSub},
      # Start the Finch HTTP client for sending emails
      {Finch, name: Bcit.Finch},
      # Start a worker by calling: Bcit.Worker.start_link(arg)
      # {Bcit.Worker, arg},
      # Start to serve requests, typically the last entry
      BcitWeb.Endpoint
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Bcit.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    BcitWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end
