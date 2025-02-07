defmodule Lab4b do
  def atomic_weight(symbol) do
    function_name = String.downcase(symbol) |> String.to_atom()

    try do
      apply(Lab4aFunctions, function_name, [])
    rescue
      UndefinedFunctionError ->
        -1
    end
  end
end
