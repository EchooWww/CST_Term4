defmodule Bcit.SchoolFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Bcit.School` context.
  """

  @doc """
  Generate a unique student sid.
  """
  def unique_student_sid, do: "some sid#{System.unique_integer([:positive])}"

  @doc """
  Generate a student.
  """
  def student_fixture(attrs \\ %{}) do
    {:ok, student} =
      attrs
      |> Enum.into(%{
        firstname: "some firstname",
        last: "some last",
        score: 42,
        sid: unique_student_sid()
      })
      |> Bcit.School.create_student()

    student
  end
end
