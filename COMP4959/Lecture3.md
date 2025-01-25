## Lecture 3

### Enum Type

- Enum in rust can be declared as follows:

```rust
enum Shape{
    Circle(f32, f32, f32),
    Square(f32),
}
```

In this example, each enum is a tuple struct with a name.

## Generics

We have a sum for i32:

```rust
fn sum (s:&[i32]) -> i32 {
    let mut sum = 0;
    for i in s {
        sum += i; // we can either dereference i or use *i, as it can be dereferenced automatically
    }
    sum
}
```

How would we write a function to sum all types of numbers?

```rust
fn sum<T>(s:&[T]) -> T {
    let mut sum = s[0];
    for x in s.into_iter().skip(1) {
        sum += *x;
    }
    sum
}

fn main () {
    let a = [1, 2, 3, 4, 5];
    println!("{}", sum(&a));
    let b = [1.0, 2.0, 3.0, 4.0, 5.0];
    println!("{}", sum(&b));
}
```

This will cause an error, as the + operator is not defined for type T. We can fix this by adding a trait `AddAssign` to the function signature:

Still, the `+=` operator requires the value to be copied. We can do this by adding the `Copy` trait to the function signature:

```rust
fn sum<T: std::ops::AddAssign + Copy>(s:&[T]) -> T {
    let mut sum = s[0];
    for x in s.into_iter().skip(1) {
        sum += *x;
    }
    sum
}
```

We can also have a generic function that takes a function as an argument:

```rust
fn find<>(s: &[T], f: fn(&T) -> bool) -> Option<&T> {
    for x in s {
      if f(x) {
        return Some(x);
      }
    }
    None
}
```

But if we call it in main() like this:

```rust
    let n = 3
    fn is_divisible(x: &i32) -> bool {
        x % n == 0
    }
```

It cannot capture the variable `n` inside the function. We can fix this by using a closure:

```rust
    let n = 3
    let is_divisible = |x: &i32| -> bool {
        x % n == 0
    }
```

Now, the anonymous function can capture the variable `n` outside of its scope.

But we cannot pass it to the function `find` as it is expecting a function pointer, not a closure.

We can fix this by changing the function signature to accept a closure. However, closure types are not named, each closure has its unique type. To write generic functions that accept closures, we can use the `Fn` trait:

Fn<:FnMut<:FnOnce (the `<:` notation means "is a subtype of")

- FnOnce: every closure implements this trait, it can be called only once
- FnMut: it mutates the variable, but it can be called multiple times
- Fn: it does not mutate the variable, and it can be called multiple times

How a closure capture the environment? And can it be called multiple times?

```rust
fn main() {
  let v = vec![1, 2, 3, 4, 5];
  let c = || println!("{v:?}"); // borrow v; impl Fn, FnMut, FnOnce
  println!("{v:?}");
  c();
  c();
}
```

Another example

```rust
  let mut v = vec![1, 2, 3, 4, 5];
  let mut c = || v.push(-1); // mutably borrow v, impl FnMut, FnOnce
  // println!("{v:?}"); // cannot borrow v again
  c();
  c(); // can call multiple times

```

Now let's take a look at FnOnce:

```rust
  let mut v = vec![1, 2, 3, 4, 5];
  let mut c = || v.push(-1); // mutably borrow v, impl ONLY FnOnce
  // println!("{v:?}"); // cannot borrow v again
  c();
  // c(); // cannot call multiple times
```

Basically, if a function takes a closure that implements `FnOnce`, we can pass any closure.

The `where` keyword can be used to specify the trait bounds for a generic type. It is equivalent to type constraints in the function signature.

```rust
fn find<T, F>(s: &[T], f: F) -> Option<&T>
  where F: Fn(&T) -> bool {
  for x in s {
    if f(x) {
      return Some(x);
    }
  }
  None
}
```

For functions, every function implements `FnOnce`, `FnMut`, and `Fn`. So, we can use `Fn` as a trait bound for functions.

So when we write generic functions accepting functions, it's better to make it accept `Fn` trait, as it could accept function pointers and closures.

## Creating an iterator

We have 3 functions can be used to create an iterator:

- iter(): borrows the collection and returns an iterator. It takes `&self` as a parameter.
- into_iter(): consumes the collection and returns an iterator. Takes ownership of the collection. It takes `self` as a parameter.
- iter_mut(): mutably borrows the collection and returns an iterator. It takes `&mut self` as a parameter.

The for loop is a syntactic sugar for iterators. It calls the `into_iter()` method on the collection, which returns an iterator. Then, it calls the `next()` method on the iterator until it returns `None`.

The IntoIterator trait has a method called into_iter() that returns an iterator, too.

How can we implement an iterator for vectors? First, the iterator must be a struct that implements the Iterator trait.
