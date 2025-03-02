defmodule Bcit.SchoolTest do
  use Bcit.DataCase

  alias Bcit.School

  describe "students" do
    alias Bcit.School.Student

    import Bcit.SchoolFixtures

    @invalid_attrs %{last: nil, sid: nil, firstname: nil, score: nil}

    test "list_students/0 returns all students" do
      student = student_fixture()
      assert School.list_students() == [student]
    end

    test "get_student!/1 returns the student with given id" do
      student = student_fixture()
      assert School.get_student!(student.id) == student
    end

    test "create_student/1 with valid data creates a student" do
      valid_attrs = %{last: "some last", sid: "some sid", firstname: "some firstname", score: 42}

      assert {:ok, %Student{} = student} = School.create_student(valid_attrs)
      assert student.last == "some last"
      assert student.sid == "some sid"
      assert student.firstname == "some firstname"
      assert student.score == 42
    end

    test "create_student/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = School.create_student(@invalid_attrs)
    end

    test "update_student/2 with valid data updates the student" do
      student = student_fixture()
      update_attrs = %{last: "some updated last", sid: "some updated sid", firstname: "some updated firstname", score: 43}

      assert {:ok, %Student{} = student} = School.update_student(student, update_attrs)
      assert student.last == "some updated last"
      assert student.sid == "some updated sid"
      assert student.firstname == "some updated firstname"
      assert student.score == 43
    end

    test "update_student/2 with invalid data returns error changeset" do
      student = student_fixture()
      assert {:error, %Ecto.Changeset{}} = School.update_student(student, @invalid_attrs)
      assert student == School.get_student!(student.id)
    end

    test "delete_student/1 deletes the student" do
      student = student_fixture()
      assert {:ok, %Student{}} = School.delete_student(student)
      assert_raise Ecto.NoResultsError, fn -> School.get_student!(student.id) end
    end

    test "change_student/1 returns a student changeset" do
      student = student_fixture()
      assert %Ecto.Changeset{} = School.change_student(student)
    end
  end
end
