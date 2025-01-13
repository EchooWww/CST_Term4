defmodule Primes do
  def primes(n) do
    if n < 2 do
      []
    else
      do_seive(Enum.to_list(2..n), [], n)
    end
  end

  defp do_seive([h | t], acc, n) do
    if h * h > n do
      Enum.reverse([h | acc]) ++ t
    else
      filtered = Enum.filter(t, fn x -> rem(x, h) != 0 end)
      do_seive(filtered, [h | acc], n)
    end
  end

  def group_permutations(primes) do
    primes
    |> Enum.group_by(&sort_digits/1)
    |> Enum.max_by(fn {_, v} -> length(v) end)
    |> elem(1)
  end

  defp sort_digits(number) do
    number
    |> Integer.to_string()
    |> String.graphemes()
    |> Enum.sort()
    |> Enum.join()
  end
end

IO.inspect(length(Primes.group_permutations(Primes.primes(1000000))))
