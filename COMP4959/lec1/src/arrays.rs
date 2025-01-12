pub fn run(){
  println!("ARRAYS"); 

  let a = [3,2,7,6,8];
  print(&a);
  print2(&a);
  let mut a = [1,2,3];
  a[0] = -1;
}

fn print(a:&[i32;5]) {
  for i in 0..5 {
    println!("{}", a[i]);
  }
}

fn print2(a:&[i32;5]) {
  for x in a {
    println!("{}", x);
  }
}

fn double(a: &mut [i32;5]) {
  for x in a { // x has type &mut i32
    // the for loop is creating an iterator, which is immutable, so we need to dereference x to make it mutable
    *x *= 2;
  }
}