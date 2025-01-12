defmodule Lec1 do
  def fact(n) do
    if n <= 0 do 1 else n * fact(n - 1) end
  end

  defp fact2(0, acc), do: acc
  defp fact2(n, acc), do: fact2(n - 1, n * acc)
  def fact2(n), do: fact2(n, 1)

  defp fact3(n, acc) when n <= 0, do: acc
  defp fact3(n, acc), do: fact3(n - 1, n * acc)
  def fact3(n), do: fact3(n, 1)

  def fact4(n) do
    case n do
      0 -> 1
      _ -> n * fact4(n - 1)
    end
  end

  def dedup(l) do
    case l do
      [h | x = [h | _]] -> dedup(x)
      [h | x = [_ | _]]-> [h | dedup(x)]
      _ -> l
    end
  end

  def fizzbuzz(n), do: fizzbuzz(1, n)

  defp fizzbuzz(i, n) do
    if i > n do
      :ok
    else
      print(i)
      fizzbuzz(i + 1, n)
    end
  end

  defp print(n) do
    cond do
      rem(n, 15) == 0 -> IO.puts("fizzbuzz")
      rem(n, 5) == 0 -> IO.puts("buzz")
      rem(n, 3) == 0 -> IO.puts("fizz")
      true -> IO.puts(n)
    end
  end
end

defmodule Shape do
  @pi 3.14159265   # module attribute
  def area(shape) do
    case shape do
      {:square, side} -> side * side
      {:circle, radius} -> @pi * radius * radius
      _ -> {:error, "i don't haw the hack to calculate the area"}
    end
  end
end
