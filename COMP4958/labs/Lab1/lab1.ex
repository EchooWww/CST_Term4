defmodule Lab1 do
  def mod_inverse(a,n) do
    find_mod_inverse(0,1,n,a,n)
  end

  defp find_mod_inverse(t, newt, r, 0, n) do
    cond do
      r > 1 -> :not_invertible
      t < 0 -> t + n
      true -> t
    end
  end

  defp find_mod_inverse(t, newt, r, a, n) do
    q = div(r,a)
    find_mod_inverse(newt, t - q*newt, a, rem(r,a), n)
  end


  def mod_pow(a,m,n) do
    find_mod_pow(a,m,n,1)
  end

  defp find_mod_pow(a,0,n,r) do
    r
  end

  defp find_mod_pow(a,m,n,r) do
    case rem(m,2) do
      0 -> find_mod_pow(rem(a*a,n), div(m,2), n, r)
      1 -> find_mod_pow(a, m-1, n, rem(a*r,n))
    end
  end
end
