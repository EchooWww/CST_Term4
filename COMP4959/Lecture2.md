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

   Be cautious: slicing uses byte indices, not character indices. For example:

   ```rust
   let s = "नमस्ते";
   let slice = &s[0..3]; // This works because 'न' is 3 bytes long.
   ```

   To work with characters, use iterators:

   ```rust
   let s1 = "Hello, world!";
   println!("Number of characters: {}", s1.chars().count());
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

Rust’s `Deref` trait allows `&String` to be used where `&str` is expected, as `Deref` converts `String` to `str`. This means the following works seamlessly:

```rust
fn print_str(s: &str) {
    println!("{}", s);
}

let s = String::from("hello");
print_str(&s); // Deref coercion makes this possible.
```

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

Methods can use `self` (for instances) and `&mut self` (for mutable references). Rust automatically dereferences or references structs when invoking methods.

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

### Notes

- `self` must always be the first parameter in methods.
- Methods expecting mutable references require mutable instances.

```rust
let mut c = Circle::new(Point(0.0, 0.0), 1.0);
c.translate(Vector(1.0, 2.0));
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
let result = square_root(-2.0).expect("Cannot take square root of a negative number");
```

Use the `?` operator for concise error handling:

```rust
fn f(x: f32) -> Option<f32> {
    Some(square_root(x)?.ln())
}
```

Note that `?` can only be used in functions returning `Result` or `Option`, because when an error occurs from `?`, the function returns early with the error.

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

Store lines in a vector:

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
println!("{:?}", chars); // ['न', 'म', 'स', '्', 'त', 'े']
```

Note: `chars()` iterates over Unicode scalar values, returning characters rather than bytes.

### Three Ways to Create Iterators

1. **`iter()`**:

   - Creates an iterator that **borrows** each element of a collection (e.g., array, vector).
   - It doesn't consume the collection, so the collection can still be used after iteration.
   - Example:
     ```rust
     let a = vec![1, 2, 3];
     let v: Vec<_> = a.iter().collect();
     assert_eq!(v, vec![&1, &2, &3]); // Note: Elements are references.
     ```

2. **`into_iter()`**:

   - Consumes the collection and creates an iterator that takes ownership of its elements.
   - After using `into_iter()`, the original collection is no longer available.
   - Example:
     ```rust
     let a = vec![1, 2, 3];
     let v: Vec<_> = a.into_iter().collect();
     assert_eq!(v, vec![1, 2, 3]); // Elements are owned, no references.
     ```

3. **`iter_mut()`**:
   - Creates a **mutable iterator** that allows modifying elements of the collection during iteration.
   - The collection must be mutable.
   - Example:
     ```rust
     let mut a = vec![1, 2, 3];
     for x in a.iter_mut() {
         *x *= 2; // Modify elements in place.
     }
     assert_eq!(a, vec![2, 4, 6]);
     ```
