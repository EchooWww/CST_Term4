defmodule Bcit.Repo.Migrations.CreateStudents do
  use Ecto.Migration

  def change do
    create table(:students) do
      add :sid, :string
      add :firstname, :string
      add :last, :string
      add :score, :integer

      timestamps(type: :utc_datetime)
    end

    create unique_index(:students, [:sid])
  end
end
