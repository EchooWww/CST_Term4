// fn sum (s:&[i32]) -> i32 {
//   let mut sum = 0;
//   for i in s {
//       sum += i;
//   }
//   sum
// }
use std::ops::AddAssign;


fn sum<T: AddAssign + Copy>(s:&[T]) -> T {
    let mut sum = s[0];
    for x in s.into_iter().skip(1) {
        sum += *x;
    }
    sum
}

fn find<T>(s: &[T], f: fn(&T) -> bool) -> Option<&T> {
  for x in s {
    if f(x) {
      return Some(x);
    }
  }
  None
}

fn find2<T, F>(s: &[T], f: F) -> Option<&T>
  where F: Fn(&T) -> bool {
  for x in s {
    if f(x) {
      return Some(x);
    }
  }
  None
}

fn main () {
    let a = [1, 2, 3, 4, 5];
    println!("{}", sum(&a));
    let b = [1.0, 2.0, 3.0, 4.0, 5.0];
    println!("{}", sum(&b));
    let c = [1,2,3,4,5];
    fn is_even(x: &i32) -> bool {
        x % 2 == 0
    }
    if let Some(x) = find(&c, is_even) {
        println!("{}", x);
    } else {
        println!("Not found");
    }

    let n = 3;
    // fn is_divisible(x: &i32) -> bool {
    //     x % n == 0
    // }
    let is_divisible = |x: &i32| x % n == 0;
    // let _ = find(&c, is_divisible);
    let _ = find2(&c, is_divisible);
}

