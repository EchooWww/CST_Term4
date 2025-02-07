defmodule Testing do
  defmacro check(pred) do
    name = Macro.to_string(pred)
    quote do
      @tests [unquote(name) | @tests]
      IO.write("Checking #{unquote(name)}: ")
      if unquote(pred) do
        IO.puts("PASSED")
      else
        IO.puts("FAILED")
      end
    end
  end

  defmacro test(name, block) do
    fn_name = String.to_atom("test_" <> name)
    quote do
      def unquote(fn_name)() do
        unquote(block)
      end
    end
  end

  defmacro __using__(_opt) do
    quote do
      import unquote(__MODULE__)
      @tests [] # module attribute
      @before_compile unquote(__MODULE__)
    end
  end

  def __before_compile__(_env) do
    quote do
      def run() do
        Enum.each(
          Enum.reverse(@tests),
          fn test -> IO.puts("Running test: #{test}"); apply(__MODULE__, test, []) end
        )
      end
    end
  end
end

defmodule Testing.Examples do
  # import Testing
  use Testing # requires a macro named __using__/1, and when `use` is called, it calls __using__/1
  test "equality" do
    check(1+1 == 2)
    check(1+1 == 3)
  end

  test "less_than" do
    check(1 < 2)
    check(2 < 1)
  end
end
