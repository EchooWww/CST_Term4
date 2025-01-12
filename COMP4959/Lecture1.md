# Lecture 1: Introduction to Rust

---

## 1. `main` Function

- The `main.rs` file is the program's entry point, created when we run `cargo new lec1`.
- It is located in the `src` directory and contains the `main` function.
- We can use the `mod` keyword to include other modules:

```rust
mod numbers;

fn main() {
    numbers::print_numbers(); // Call the print_numbers function from the numbers module
}
```

---

## 2. Type Inference

- Rust infers variable types automatically when declared using `let`:

```rust
let x = 5; // x is inferred to be an integer
```

- Rust looks at the whole file to resolve types. Any conflicting types will result in a compile-time error:

```rust
let x = 5;
println!("{}", x + 256); // Error: `256` does not fit into the range of `u8` (0..=255)

println!("{}", x + 1u8); // Adding `1u8` forces `x` to be inferred as a `u8`
```

---

## 3. Mutability

- Variables are immutable by default. To make them mutable, use the `mut` keyword:

```rust
let mut x = 12;

fn square(x: i32) -> i32 {
    x * x // Returns the square of x
}
```

---

## 4. Borrowing

- Borrowing allows passing a reference to a variable without transferring ownership.
- Rules for borrowing:
  - You can have multiple **immutable references**.
  - You can have only one **mutable reference** at a time.

Example:

```rust
fn main() {
    let mut x = 5;
    let y = &x;        // Immutable reference
    let z = &mut x;    // Mutable reference
}
```

- Reassigning a variable after borrowing it will cause an error if the reference is used:

```rust
fn main() {
    let mut x = 5;
    let y = &x; // Immutable reference
    x = 6;      // Reassigning x invalidates y
    println!("{}", y); // Error: y cannot be used after x is reassigned
}
```

---

## 5. Ownership

- Assigning a variable to another variable transfers ownership, meaning the original variable becomes invalid:

```rust
let v = vec![1, 2, 3];
let v2 = v; // Ownership of v is moved to v2; v cannot be used anymore
```

---

## 6. Printing Arrays

- Arrays cannot be printed directly using `println!`. Instead, define a custom function:

```rust
fn print(a: [i32; 5]) {
    for i in 0..5 {
        println!("{}", a[i]);
    }
}
```

---

## 7. Problems with Passing Values

- Passing a value to a function transfers ownership, preventing further use of the value:

```rust
fn sum(v: Vec<i32>) -> i32 {
    v.into_iter().sum()
}

pub fn run() {
    let v = vec![3, 2, 7, 6, 8];
    println!("{:?}", v);
    println!("{}", sum(v)); // Ownership of v is consumed
    // println!("{:?}", v); // Error: v is no longer valid
}
```

- Use references to avoid consuming the value:

```rust
fn sum(v: &Vec<i32>) -> i32 {
    v.iter().sum()
}

pub fn run() {
    let v = vec![3, 2, 7, 6, 8];
    println!("{:?}", v);
    println!("{}", sum(&v)); // Passing a reference allows reuse of v
    println!("{:?}", v);     // No error
}
```

---

## 8. Arrays, Vectors, and Slices

- **Arrays**: Fixed size, stack-allocated.
- **Vectors**: Dynamically sized, heap-allocated.
- **Slices**: References to arrays or vectors, providing a view of the data.

Example:

```rust
fn double(a: &mut [i32]) {
    for x in a {
        *x *= 2; // Dereference mutable reference to modify values
    }
}
```

---

## 9. Visibility

- By default, everything in Rust is private. Use the `pub` keyword for public visibility.
- To call functions from another module:

```rust
let mut v = vec![1, 2, 3];
super::double(&mut v); // Access the double function from a parent module
```

---

## 10. Higher-Order Functions

- Functions that take or return other functions are called higher-order functions.

Example:

```rust
pub fn find(s: &[i32], f: fn(i32) -> bool) -> Option<i32> {
    for &x in s {
        if f(x) {
            return Some(x);
        }
    }
    None
}

fn is_even(x: i32) -> bool {
    x % 2 == 0
}

pub fn run() {
    let v = vec![1, 2, 3, 4, 5];
    if let Some(x) = find(&v, is_even) {
        println!("First even number is {}", x);
    }
}
```

---

## 11. `if let`

- `if let` is a concise way to match a single pattern:

```rust
if let Some(x) = find(&v, is_even) {
    println!("First even number is {}", x);
} else {
    println!("No even number found");
}
```

---
