pub fn run(){
  println!("NUMBERS");
  let x = 1;     // rust can infer the type of x, default of integer can be i32
  // println!("{}", x + 256);
  println!("{}", x + 1u8); // x needs to be a u8: unsigned 8-bit integer. Type inference looks at the whole file to determine the type of x, If there are conflicting types, rust will throw an error.
  let mut x = 12;
  println!("{}", square2(&x));
  triple(&mut x);
  println!("{}", x);
  println!("{}", &x); // still printing the value
  println!("{}", &&x); // still printing the value
  let r = &x;
  // x = 1; // no lexical lifetime, so x can be reassigned
  println!("{}", r); // but if we try to use r, it will throw an error: cannot borrow `x` as mutable because it is also borrowed as immutable
  let x = 14;
  println!("{}", x + &1);
}

fn square (x:i32) -> i32 {
  x * x
}

fn square2(x:&i32) -> i32 {
  x * x
}

fn triple(x:&i32)-> i32 {
  *x * 3
}
