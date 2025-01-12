# Lecture 1: Intro to Elixir

## 1. Interactive Shell (`iex`)

- `h <module>`: Get help on a module.
- `h <module>.<function>`: Get help on a function.
  - Use `Tab` for autocompletion of available functions.

---

## 2. Data Types in Elixir

Elixir is functional and dynamically typed, with inferred types.

### 2.1 Enum Module

- Provides a set of functions to work with enumerables.
- Lists are enumerables with additional specific functions.

### 2.2 Atoms

- Atoms are constants where their name is their value (e.g., `:ok`, `:error`, `:nil`).
- Often used to represent states or tags.

Examples:

```elixir
{:ok, "Hello"} # Tuple
{:error, "Something went wrong"}
```

- Atoms with spaces can be wrapped in `:""`.

```elixir
:"Hello, World"
```

- Booleans (`true` and `false`) are atoms.

```elixir
true == :true # Returns true
```

- Anything starting with an uppercase letter is also an atom.

```elixir
:Hello
```

### 2.3 Lists

- Immutable and can contain multiple data types.

```elixir
l = [1, 2, 3, 4, 5, "Hello", :ok]
```

---

## 3. Operators

### 3.1 Equal Sign (`=`)

- In Elixir, `=` is a match operator, not an assignment operator.
- Matches the right side with the left side.
  - If unbound, binds the value.
  - Otherwise, checks for a match.

Examples:

```elixir
x = 10
20 = x # Throws an error as 20 does not match x (10)
```

### 3.2 Pin Operator (`^`)

- Prevents rebinding of a variable.

```elixir
x = 10
^x = 20 # Throws a match error instead of reassigning
```

### 3.3 Cons Operator (`|`)

- Used to build lists or pattern match.

```elixir
[1 | [2, 3, 4]] # Equivalent to [1, 2, 3, 4]
```

Example with pattern matching:

```elixir
[_,_,{_,x}|_] = [1, 2, {:ok, "Hello"}, 4, 5]
```

---

## 4. Modules and Functions

### 4.1 Defining Modules

```elixir
defmodule M do
  def fact(n) do
    if n <= 0 do 1 else n * fact(n-1) end
  end
end
```

### 4.2 Tail Recursion

```elixir
def fact2(0, acc), do: acc
def fact2(n, acc), do: fact2(n-1, n*acc)
def fact2(n), do: fact2(n, 1)
```

Private functions:

```elixir
defmodule M do
  def fact(n), do: fact2(n, 1)

  defp fact2(0, acc), do: acc
  defp fact2(n, acc), do: fact2(n-1, n*acc)
end
```

### 4.3 Using Guards

```elixir
defmodule M do
  def fact(n), do: fact3(n, 1)

  defp fact3(0, acc), do: acc
  defp fact3(n, acc) when n > 0, do: fact2(n-1, n*acc)
end
```

---

## 5. Input/Output

### 5.1 Printing to Console

```elixir
IO.puts "Hello, World"
# Prints "Hello, World" and returns :ok
```

---

## 6. Control Flow

### 6.1 Using `cond`

- Checks multiple conditions; the first that matches is executed.

```elixir
defmodule Lec1 do
  def fizzbuzz(n), do: fizzbuzz(1, n)
  defp fizzbuzz(i, n) do
    if i > n do
      :ok
    else
      print(i)
      fizzbuzz(i+1, n)
    end
  end

  defp print(n) do
    cond do
      rem(n, 15) == 0 -> IO.puts "FizzBuzz"
      rem(n, 5) == 0 -> IO.puts "Buzz"
      rem(n, 3) == 0 -> IO.puts "Fizz"
      true -> IO.puts n
    end
  end
end
```

---

## 7. Anonymous Functions

```elixir
f = fn x -> x * 2 end
f.(2) # Call with . operator
```

Using as parameters:

```elixir
Enum.map(["hello", "world"], fn x -> String.upcase(x) end)
# Or with the capture operator:
Enum.map(["hello", "world"], &String.upcase/1)
```

---

## 8. Pipe Operator (`|>`)

- Passes the result of the left side as the first argument to the right.

```elixir
[1, 2, 3] |> Enum.map(&(&1 * &1)) |> Enum.sum
```

---

## 9. Compilation

### Using `elixirc`

- Compile a `.ex` file to generate a `.beam` file.

Example:

```elixir
defmodule Shape do
  @pi = 3.14159265
  def area(shape) do
    case shape do
      {:square, side} -> side * side
      {:circle, radius} -> @pi * radius * radius
      _ -> {:error, "Invalid shape"}
    end
  end
end
```
