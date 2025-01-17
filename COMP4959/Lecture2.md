## Strings

There are 2 types of strings in Rust:
`let s1 = "Hello, world!";` // String slice, it's a reference to a string in the binary

We cannot index into a string slice:

```rust
let s1 = "Hello, world!";
let h = s1[0]; // Error: cannot index into a string slice
```

`rustc` is the command-line compiler for Rust. It compiles Rust code into an executable binary.

The indexing is related to a trait. Only types that implement the `Index` trait can be indexed into. The `Index` trait is not implemented for string slices.

Instead, we can slice the string:

```rust
let s1 = "Hello, world!";
let h = &s1[0..1]; // Slices the string from index 0 to 1, excluding 1
```

But taking slices is also dangerous: it is using the index of a byte, not a character. This is because Rust strings are UTF-8 encoded. Each character is a codepoint, a codepoint could be multiple bytes. `len()` returns the number of bytes, not the number of characters.

To get the characters, we get an iterator:

```rust

let s1 = "Hello, world!";
println("{}", s1.chars().count());
```

Another type of string is created with `String::from()`. It is a String (Object) type, not a string slice.

If we wanna pass a String to a function expecting a string slice, it won't compile. We need to pass a reference to the String:

```rust
fn main() {
  let s2 = String::from("hello");
  println!("{} {} {}", s, s.len(), s.capacity());
  print_str(&s2);
  print_string_ref(&s2);
  print_string(s2);
}

fn print_str(s: &str) {
  println!("{}", s);
}

fn print_string(s: String) {
  println!("{}", s);
}

fn print_string_ref(s: &String) {
  println!("{}", s);
}

```

- That means, when we wanna use a string, it's better to use the String object, as it can be passed to functions expecting both `&str` and `String`.

- Likewise, when writing functions, it's better to use `&str` as the parameter type, as it can accept both `&str` and `&String`.

> However, if we move this line `print_string(s2);` before `print_string_ref(&s2);`, it won't compile. This is because the `print_string` function takes ownership of the string, and we can't use it anymore. In this case, we need to clone the string:
>
> ```rust
> fn main() {
>   let s2 = String::from("hello");
> println!("{} {} {}", s, s.len(), s.capacity());
>   print_str(&s2);
>   print_string(s2.clone());
>   print_string_ref(&s2);
> }
> ```

There is a `Deref` trait that allows us to dereference a type. That is why we can use `&String` where a `&str` is expected: the `Deref` trait is implemented for `String`, and the target type is `str`.

`Deref` takes a reference to itself, and return a reference to the target type. It is a way to convert a type to another type.

## Structs

In Rust, structs are similar to classes in other languages. They are used to create custom data types.

We can also implement methods for a struct with the `impl` keyword. Inside the implementation, we can use `Self` to refer to the struct type. There's also a lower case `self` that refers to an instance of the struct.

```rust
struct Point(f32, f32);
struct Vector(f32, f32);
struct Circle {
    center: Point,
    radius: f32,
}

impl Circle {
  fn new(center: Point, radius: f32) -> Self {
    Self { center, radius }
  }
}
```

Using `self`. Note we're using `&mut self` to borrow a mutable reference to the struct. This is because we're modifying the struct.

```rust
  fn scale(&mut self, factor: f32) {
    self.radius *= factor;
  }
```

We can "translate" a struct by adding a vector to its center. Note we access elements in a tuple using `.`:

```rust
  fn translate(&mut self, v: Vector) {
    self.center.0 += v.0;
    self.center.1 += v.1;
  }
```

We we invoke a method on a struct, Rust automatically reference or dereference the struct as needed, no matter how many times.

```rust
  fn main() {
  let c = Circle::new(Point(0.0, 0.0), 1.0);
  println!("Circle area: {}", c.area());
}
```

Also, `self` can only be the first parameter of a method.

We can also call that method as a function, but then we need to explicitly pass the reference:

```rust
  fn main() {
  let c = Circle::new(Point(0.0, 0.0), 1.0);
  println!("Circle area: {}", Circle::area(&c));
}
```

To invoke methods expecting a mutable reference, we need a mutable object, and it will be a mutable borrow automatically.

```rust
  fn main() {
  let mut c = Circle::new(Point(0.0, 0.0), 1.0);
  c.scale(2.0);
  println!("Circle area: {}", c.area());
}
```

Now we want to print a Circle: to do that, we need to implement the `Display` trait for the Circle struct. The `Display` trait is used to format a value using the `{}` placeholder.

We can use this command `rustup doc std::fmt::Display` to see the documentation for the `Display` trait.

Inside the `fmt` function,

```rust
use std::fmt;
impl fmt::Display for Circle {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    write!(f, "Circle with center at ({}, {}) and radius {}", self.center.0, self.center.1, self.radius)
  }
}
```

Note `Result` itself is an enum, with variants `Ok` and `Err`. Other structs can implement Result, so we have `fmt::Result` to return the specific variant of Result.

## Error Handling

```rust
fn square_root(x:f32) -> Option<f32> {
  if x < 0.0 {
    return None;
  }
  Some(x.sqrt())
}

fn main() {
  let x = square_root(2.0).unwrap(); // unwrap() may panic if the value is None
  let x = square_root(2.0).expect("Cannot take square root of a negative number"); // expect() allows us to specify a custom error message
}
```

A more sophisticated way to handle:

```rust
fn f (x:f32) -> Option<f32> {
  Some(square_root(x)?.ln()) // The ? operator returns the value if it's Some, otherwise returns None
}
```

The question mark operator is a shorthand for the `match` statement. It returns the value if it's `Some`, otherwise returns `None` instantly from the function.

## IO

The difficulty of IO is that there are so many ways to do it. Rust has a standard library that provides a lot of IO functionality: Read, BufRead, BufReader...

### Read from stdin

`io::stdin().read_line()` reads a line from stdin and returns a `Result` with the number of bytes read: either `Ok(n)` or `Err(e)`. When n is 0, it means EOF.

```rust
use std::io;
fn main() {
  let stdin = io::stdin();
  let mut line = String::new();
  while let Ok(n) = stdin.read_line(&mut line) {
    if n ==0 {
      break;
    }
    print!("{line}");
    line.clear(); // clear the string for the next line
  }
}
```

### Read multiple lines from stdin

Read lines and store them in a vector. Note that when pushing a string to a vector, we need to clone it. If we push it directly, the String will be consumed and we can't use it anymore for the next iteration.

```rust
fn read_lines() -> io::Result<Vec<String>> {
  let mut line = String::new();
  let mut v = Vec::new();
  while io::stdin().read_line(&mut line)? > 0 {
    v.push(line.clone());
    line.clear();
  }
  Ok(v)
}
```

With the use of `?`, we can return the error from the function. If the `read_line` function returns an error, the error will be returned from the function.

### Read from a file

The only difference is that we need to open the file first. We can use the `File` struct from the `std::fs` module.

```rust

fn read_file_lines(file: &str) -> io::Result<Vec<String>> {
  let f = File::open(file)?;
  let mut reader = BufReader::new(f);
  let mut line = String::new();
  let mut v = Vec::new();
  while reader.read_line(&mut line)? > 0 {
    v.push(line.clone());
    line.clear();
  }
  Ok(v)
}
```

## Collect

Collect is a method that converts an iterator into a collection. It is used to collect the results of an iterator into a collection.

3 ways to create an interator:

- `iter()`: creates an iterator over a container
- `into_iter()`: consumes the container and creates an iterator
- `iter_mut()`: creates a mutable iterator over a container

To use `into_iter()`, we need to run `rustc --edition 2021 collect.rs` to specify the edition.

```rust
let a = vec![1,2,3];
let v:Vec<_> = a.into_iter().collect();
  println!("{:?}", v);
  assert_eq!(v, vec![1,2,3]);
```

Note we can specify the type of the collection returned by `collect`. In this case, we're using `Vec<_>` to let Rust infer the type.

`iter()`, on the other hand, borrows each element of the collection. This is useful when we don't want to consume the collection.

```rust
let a = vec![1,2,3];
let v:Vec<_> = a.iter().collect();
  assert_eq!(v, vec![&1,&2,&3]);
```

Making use of the type inference, we can convert an array of key-value pairs into a HashMap:

```rust
let a = [("a", 1), ("b", 2), ("c", 3)];
let m:HashMap<_,_> = a.iter().collect();
```

As strings cannot be indexed, it's a good idea to convert a string to a vector of characters using `chars()` and `collect()`:

```rust
let s = "नमस्ते";
let v:Vec<_> = s.chars().collect();
```

It will return a vector of characters, instead of a vector of bytes.
