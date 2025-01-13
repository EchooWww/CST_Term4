defmodule ArithmeticServer do
  def start() do
    spawn(&loop/0)
  end

  def square(pid, x) do
    send(pid, {:square, x, self()})
    receive do
      x -> x
    end
  end

  def sqrt(pid, x) do
    send(pid, {:sqrt, x, self()})
    receive do
      x->x
    end
  end

  defp loop() do
    receive do
      {:square, x, from} ->
        send(from, x * x)
      {:sqrt, x, from} ->
	response = 
	  if x >= 0, do: {:ok, :math.sqrt(x)}, else: {:error, "Negative number"}
        send(from, response)
    end
    loop()
  end
end

def CounterServer do
  def start(n\\0) do
    spawn(fn -> loop(n) end)
  end
  
  def inc(pid) do
    send(pid, :inc)
  end

  def dec(pid) do
    send(pid, :dec)
  end

  def value(pid) do
    send(pid, {:value, self()})
  end

  defp loop(n) do
    receive do
      :inc ->
        loop(n+1)
      :dec ->
        loop(n-1)
      {:value, from} ->
        send(from, n)
	loop(n)
      end

end
    end
