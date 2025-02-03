enum Shape{
    Circle(f32, f32, f32),
    Square(f32),
}

impl Shape{
  fn area(&self) -> f32 {
    match self {
      Shape::Circle(_, _, r) => std::f32::consts::PI * r * r,
      Shape::Square(s) => s * s,
    }
  }
}

fn sum (s:&[i32]) -> i32 {
  let mut sum = 0;
  for i in s {
      sum += i;
  }
  sum
}

fn main () {
  let a = [Shape::Circle(1.0, 1.0, 2.0), Shape::Square(2.0)];
  a.iter().for_each(|s| println!("{:.2}", s.area()));

  let b = [1, 2, 3, 4, 5];
  println!("{}", sum(&b));
}