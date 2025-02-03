pub fn run() {
  println!("VECTORS");
  let v = vec![3,2,7,6,8];
  println!("{v:?}");
  // println!("{}", sum(v)); // the side effect of sum is that it consumes the vector
  // println!("{v:?}"); // so v is no longer available
  println!("{}", super::slices::sum(&v)); // We can only call the function in a lower module if it is public
  let mut v = vec![1,2,3];
  super::double(&mut v);// we can call the function in a higher module no matter if it is public or private
  println!("{v:?}"); 
}

fn sum(v:Vec<i32>) -> i32 {
  let mut s = 0;
  for x in v {
    s += x;
  }
  s
}