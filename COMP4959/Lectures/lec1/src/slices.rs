pub fn run() {
  println!("SLICES");
  let s = &[1,2,3,4,5];
  println!("{:?}", s);
  let v = vec![1,2,3,4,5];
  fn is_even(x:i32) -> bool {
    x % 2 == 0
  }
  if let Some(x) = find(&v, is_even) {
    println!("First even number is {}", x);
  } else {
    println!("No even number found");
  }
}

pub fn sum(s:&[i32]) -> i32 {
  let mut sum = 0;
  for x in s {
    sum += x;
  }
  sum
}

pub fn find(s: &[i32], f:fn(i32) -> bool) -> Option<i32> {
  for x in s {
    if f(*x) {
      return Some(*x);
    }
  }
  None
}