# Lecture 5: Meta Programming with Elixir

We wanna implement a custom "if" in a function, but is it possible to do that in Elixir? No, because Elixir is not a lazy language, so it will evaluate all the arguments before calling the function.

That's when we need meta programming. We can use macros to achieve that. Macros are functions that generate code. They are executed at compile time, not runtime.

## AST: Abstract Syntax Tree

When the compiler reads the code, there are always some steps:

- Lexing: Convert the code into tokens
- Parsing: Convert the tokens into an AST. Macros are used to manipulate the AST.

Quote is a macro that converts the code into an AST. `quote do: 1 + 2` will return `{:`, [], [{:+, [context: Elixir, import: Kernel], [1, 2]}]}`.

The return value is a triple: the module, the metadata, and the operands.

There's an `unquote` function that allows us to inject values of quoted expressions.

```elixir
x = quote do: 1 + 2
quote do: 3 * unquote(x)
```

What quote do is: tell the complier not to evaluate the expression, but to return the AST. Unquote says "evaluate this expression and put it here". If the `unquote` gives a valid AST, then the whole thing will be a valid AST. If not, it will throw an error.

There's a `Macro.to_string` function that converts an AST back to a string.

```elixir
iex(1)> x = quote do: 1 + 1
{:+, [context: Elixir, imports: [{1, Kernel}, {2, Kernel}]], [1, 1]}
iex(2)> quote do: x
{:x, [], Elixir}
iex(3)> Macro.escape(x)
{:{}, [], [:+, [context: Elixir, imports: [{1, Kernel}, {2, Kernel}]], [1, 1]]}
```

`Macro.escape` will escape the AST, so it will evaluate X and put it in the AST.

When we define a macro as `defmacro foo(x) do`, and we call it as `foo(1 + 2)`, the `x` will be the AST of `1 + 2`, instead of 3.

A simpliest example of a macro is:

```elixir
IO.inspect("Hello") # print something and return the value
```

A macro takes an AST(we pass in an expression, and it will be converted to an AST), and returns an AST.

So we're gonna have a macro that does the same thing:

```elixir
defmodule MyMacro do
  defmacro my_inspect(x) do
      IO.inspect(x)
  end
end
```

AST literals: for numbers, atoms, strings, lists, pairs.. (anything that is not a variable), the AST is the value itself.

Now we can have our own `if` with a macro:

```elixir
defmacro if_macro(predicate, do_block, else_block) do
  quote do
    cond do
      unquote(predicate) -> unquote(do_block)
      true -> unquote(else_block)
    end
  end
end
```

To call a module with a macro, we need to use `require MyMacro` and `MyMacro.my_inspect("Hello")`.

We can have a better version of pattern matching with macros:

```elixir
defmacro if(predicate, do: do_block, else: else_block) do
    quote do
      if unquote(predicate) do
        unquote(do_block)
      else
        unquote(else_block)
      end
    end
  end
```

Now we can call it as `if 1 == 1, do: IO.inspect("Hello"), else: IO.inspect("World")`.

## Define a while loop

```elixir
pid = spawn (fn -> Process.sleep(60000) end)
defmacro while(predicate, do: block) do
      quote do
        try do
          for _ <- Stream.cycle([:ok]) do
            if unquote(predicate), do: unquote(block), else: throw(:break)
          end
        catch
          :throw, :break -> :ok
      end
    end
  end
require Lec5
Lec5.while Process.is_alive(pid), do: Process.sleep(1000) IO.puts("Still alive")
```
