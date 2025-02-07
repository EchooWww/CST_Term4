defmodule Lab4a do
  defmacro define_elements_from_file(file_path) do
    elements =
      File.read!(file_path)
      |> String.split("\n", trim: true)
      |> Enum.map(fn line ->
        case String.split(line, ~r/\s+/, trim: true) do
          [_index, symbol, _name, weight] ->
            {String.downcase(symbol), String.to_float(weight)}

          _ ->
            nil
        end
      end)
      |> Enum.reject(&is_nil/1)

    for {symbol, weight} <- elements do
      quote do
        def unquote(String.to_atom(symbol))(), do: unquote(weight)
      end
    end
  end
end

defmodule Lab4aFunctions do
  require Lab4a

  Lab4a.define_elements_from_file("atomic-weights.txt")
end
