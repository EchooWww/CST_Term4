fn main() {
  let v = vec![1, 2, 3, 4, 5];
  let c = || println!("{v:?}"); // implements Fn
  println!("{v:?}");
  c();
  c();
  let mut v = vec![1, 2, 3, 4, 5];
  // let mut c = || v.push(-1); // implements FnMut
  println!("{v:?}");
  c();
  c();

  let s = "hello".to_string();
  let c = || s;
  // println!("{s}"); // s has moved
  c()
  // c()  // cannot be invoked twice
}