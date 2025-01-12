defmodule Lec1 do
  def fizzbuzz(n), do: fizzbuzz(1,n)
  defp fizzbuzz(i, n) do
    if i > n do
      :ok
    else
      print(i)
      fizzbuzz(i+1, n)
    end
  end

  defp print(n) do
    cond do
      rem(n,15) == 0 -> IO.puts "FizzBuzz"
      rem(n,5) == 0 -> IO.puts "Buzz"
      rem(n,3) == 0 -> IO.puts "Fizz"
      true -> IO.puts n
    end
  end
end
defmodule Shape do
  @pi=3.14159265 # module attributes
  def area(shape) do
    case shape do
      {:square, side} -> side * side
      {:circle, radius} -> @pi * radius * radius
      _ -> {:error, "Nooooo"}
    end
  end
end
