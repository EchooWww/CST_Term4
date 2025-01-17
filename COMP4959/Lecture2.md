## Strings

There are two types of strings in Rust:

1. **String slice**: A reference to a string in the binary:

   ```rust
   let s1 = "Hello, world!";
   ```

   String slices cannot be indexed directly:

   ```rust
   let s1 = "Hello, world!";
   let h = s1[0]; // Error: cannot index into a string slice
   ```

   This is because Rust strings are UTF-8 encoded, and indexing operates on bytes, not characters. To extract a portion, slicing is used:

   ```rust
   let s1 = "Hello, world!";
   let h = &s1[0..1]; // Slices the string from index 0 to 1, excluding 1
   ```

   Be cautious: slicing uses byte indices, not character indices. For character operations, use iterators:

   ```rust
   let s1 = "Hello, world!";
   println!("{}", s1.chars().count());
   ```

2. **String object**: Created with `String::from()`:

   ```rust
   let s2 = String::from("hello");
   ```

   When passing a `String` to a function expecting a string slice, pass a reference:

   ```rust
   fn print_str(s: &str) {
       println!("{}", s);
   }

   let s2 = String::from("hello");
   print_str(&s2);
   ```

   Prefer using `String` objects as they can be passed to functions expecting both `&str` and `String`. Conversely, define functions using `&str` to accept both `&str` and `&String`.

### Ownership and Cloning

If a function takes ownership of a `String`, clone it beforehand if needed elsewhere:

```rust
fn print_string(s: String) {
    println!("{}", s);
}

let s2 = String::from("hello");
print_string(s2.clone());
print_str(&s2);
```

### Deref Coercion

Rust’s `Deref` trait allows `&String` to be used where `&str` is expected, as `Deref` converts `String` to `str`.

## Structs

Structs in Rust are custom data types, similar to classes in other languages. Methods for structs are implemented using the `impl` block:

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

    fn scale(&mut self, factor: f32) {
        self.radius *= factor;
    }

    fn translate(&mut self, v: Vector) {
        self.center.0 += v.0;
        self.center.1 += v.1;
    }
}
```

Methods can use `Self` (for type) `self` (for instances) and `&mut self` (for mutable references). Rust automatically dereferences or references structs when invoking methods.

```rust
let mut c = Circle::new(Point(0.0, 0.0), 1.0);
c.scale(2.0);
```

To format a struct, implement the `Display` trait:

```rust
use std::fmt;

impl fmt::Display for Circle {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Circle with center at ({}, {}) and radius {}", self.center.0, self.center.1, self.radius)
    }
}
```

## Error Handling

Use `Option` for optional values:

```rust
fn square_root(x: f32) -> Option<f32> {
    if x < 0.0 {
        None
    } else {
        Some(x.sqrt())
    }
}

let result = square_root(2.0).unwrap(); // May panic if None
let result = square_root(2.0).expect("Cannot take square root of a negative number");
```

Use the `?` operator for concise error handling. It will return `None` if the result is `None`, and unwrap the value otherwise:

```rust
fn f(x: f32) -> Option<f32> {
    Some(square_root(x)?.ln())
}
```

## IO

Rust’s standard library provides IO utilities:

### Reading from stdin

Read lines from stdin:

```rust
use std::io;

fn main() {
    let stdin = io::stdin();
    let mut line = String::new();
    while let Ok(n) = stdin.read_line(&mut line) {
        if n == 0 {
            break;
        }
        println!("{}", line);
        line.clear();
    }
}
```

### Reading multiple lines

Store lines in a vector. Note we need to clone the line when adding it to the vector, to avoid ownership issues:

```rust
fn read_lines() -> io::Result<Vec<String>> {
    let mut line = String::new();
    let mut v = Vec::new();
    while io::stdin().read_line(&mut line)? > 0 {
        v.push(line.clone()); // Clone the line to avoid cosuming it
        line.clear();
    }
    Ok(v)
}
```

### Reading from a file

Use `File` and `BufReader`:

```rust
use std::fs::File;
use std::io::{self, BufReader, BufRead};

fn read_file_lines(file: &str) -> io::Result<Vec<String>> {
    let f = File::open(file)?;
    let reader = BufReader::new(f);
    let mut lines = Vec::new();
    for line in reader.lines() {
        lines.push(line?);
    }
    Ok(lines)
}
```

## Iterators and Collect

Iterators can be converted into collections using `collect()`:

```rust
let a = vec![1, 2, 3];
let v: Vec<_> = a.into_iter().collect();
assert_eq!(v, vec![1, 2, 3]);
```

To convert strings to characters:

```rust
let s = "नमस्ते";
let chars: Vec<_> = s.chars().collect();
```

Before collecting, we can manipulate the iterator using `map`, `filter`, `enumerate`, etc. Just like in Python:

```rust
let a = vec![1, 2, 3];
let v: Vec<_> = a.into_iter().map(|x| x * 2).collect();
assert_eq!(v, vec![2, 4, 6]);
```
