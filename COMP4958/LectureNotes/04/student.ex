defmodule Name do
  defstruct [:first, :last]

  def new(firstname, lastname) do
    %Name{first: firstname, last: lastname}
  end
end

defmodule Student do
  defstruct [id: "", name: %Name{}, scores: %{}]

  def new(id, firstname, lastname, scores \\ %{}) do
    %Student{id: id, name: Name.new(firstname, lastname), scores: scores}
  end

  def parse(data) do
    case String.split(data) do
      [id, firstname, lastname | rest] ->
        scores = 
          for [c, s] <- Enum.chunk_every(rest, 2), into: %{}, do: {c, String.to_integer(s, 10)}
        {:ok, Student.new(id, firstname, lastname, scores)}
      _ ->
        :error
    end
  end

  def read_file(path) do
    {:ok, content} = File.read(path)
    Enum.map(String.split(content, "\n"), &parse/1) |>
    Enum.filter(&(&1 != :error)) |>
    Enum.map(fn {_, x} -> x end)
  end
end
