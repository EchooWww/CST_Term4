fn square_root(x:f32) -> Option<f32> {
  if x < 0.0 {
    return None;
  }
  Some(x.sqrt())
}

fn f (x:f32) -> Option<f32> {
  // if let Some(y) = square_root(x) {
  //   Some(y.ln()) // Returns the natural logarithm of y
  // } else {
  //   None
  // }
  Some(square_root(x)?.ln()) // The ? operator returns the value if it's Some, otherwise returns None
}

fn main() {
  let x = square_root(2.0).unwrap();
  println!("Square root of 2: {}", x);
  let y = square_root(-2.0).unwrap();
  println!("Square root of -2: {}", y); // panic!
  let z = f(2.0);
  println!(z); // Some(0.6931472)
  let w = f(-2.0);
  println!(w); // None
}