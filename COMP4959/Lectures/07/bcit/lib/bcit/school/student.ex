defmodule Bcit.School.Student do
  use Ecto.Schema
  import Ecto.Changeset

  schema "students" do
    field :last, :string
    field :sid, :string
    field :firstname, :string
    field :score, :integer

    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(student, attrs) do
    attrs = for {key, value} <- attrs, into: %{}, do: {key, String.trim(value)}
    student
    |> cast(attrs, [:sid, :firstname, :last, :score])
    |> validate_required([:sid, :firstname, :last, :score])
    |> validate_format(:sid, ~r/^\d{8}$/, message: "must be 8 digits")
    |> validate_number(:score, greater_than_or_equal_to: 0, less_than_or_equal_to: 100, message: "must be between 0 and 100")
    |> unique_constraint(:sid)
  end
end
