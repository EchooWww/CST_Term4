# Lecture 4: Rust and WebAssembly

## Rust and WebAssembly

WebAssembly is a new binary instruction format for a stack-based virtual machine. It is designed as a portable target for compilation of high-level languages like C/C++/Rust, enabling deployment on the web for client and server applications.

Install wasm-pack:

```bash
brew install wasm-pack
```

Create a new Rust project:

```bash
wasm-pack new hello-wasm
cd hello-wasm
```

There's a src/lib.rs file that contains the Rust code. The code is compiled to WebAssembly using wasm-pack build.

Compile the Rust code to WebAssembly:

```bash
wasm-pack build --target web
```

Then the pkg directory contains the WebAssembly code and JavaScript bindings.

Let's implement a game of life in Rust and WebAssembly.

The pattern matching is:

```rust
_, 3 -> Alive
Alive, 2 -> Alive
_, _ -> Dead
```

In `src/lib.rs`:

```rust
use wasm_bindgen::prelude::*;
```

We need a struct to represent the Cell:

```rust
#[wasm_bindgen]
pub enum Cell {
    Dead = 0,
    Alive = 1,
}
```

The Universe struct, in which we will store the width, height, and cells:

```rust
#[wasm_bindgen]
pub struct Universe {
    width: u32,
    height: u32,
    cells: Vec<Cell>,
}
```

We cannot use a vector of vector, as we're not guarenteed that the cells are contiguous.

The implementation of the Universe struct:

```rust
#[wasm_bindgen]
impl Universe {
    fn get_index (&self, row: u32, column: u32) -> usize {
        (row * self.width + column) as usize
    }
}
```

We use usize as the return type because u32 could differ in size between the host and the WebAssembly.

#[repr(u8)] is used to represent the enum as a u8.

We also wanna add #[derive(Clone, Copy, Debug)] to enable cloning, copying, and debugging in the Cell enum.

```rust

```

We also need to implement the tick method:

````rust


To print the universe, we need to implement the Display trait:

```rust

````

In the html file, we use <pre> tag to display the universe.

```html

```

## Second version: Canvas

To make it work with the canvas, we need to add new methods to the Universe struct: width, height, cells.

```rust
    pub fn width(&self) -> u32 {
        self.width
    }
    pub fn height(&self) -> u32 {
        self.height
    }
    pub fn cells(&self) -> *const Cell {
        self.cells.as_ptr()
    }
```

We need to return a pointer to the cells, as the WebAssembly cannot return a vector.

## The trait Object

The trait Object is used to pass an object implementing a trait to a function. We need the trait Object, as the compiler cannot determine the size of the object at compile time otherwise.

```rust
fn total_area(shapes: &Vec<Box<dyn Shape>>) {
  let mut total = 0;
  for shape in shapes {
    total += shape.area();
  }
  total
}
```
