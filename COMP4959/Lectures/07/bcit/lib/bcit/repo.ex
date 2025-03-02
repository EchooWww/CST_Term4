defmodule Bcit.Repo do
  use Ecto.Repo,
    otp_app: :bcit,
    adapter: Ecto.Adapters.Postgres
end
