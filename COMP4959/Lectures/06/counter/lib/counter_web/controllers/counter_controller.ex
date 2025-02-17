defmodule CounterWeb.CounterController do
  use CounterWeb, :controller

  def value(conn, _param) do
    value = Counter.Worker.value()
    json(conn, %{value: value})
  end

  def inc(conn, %{"amt" => amt}) when is_integer(amt) do
    Counter.Worker.inc(amt)
    json(conn, :ok)
  end

  def inc(conn, %{"amt" => amt}) do
    case Integer.parse(amt) do
      {n, ""} ->
        Counter.Worker.inc(n)
        json(conn, :ok)
      _ -> 
        json(conn, %{error: "not a number"})
    end
  end
end
