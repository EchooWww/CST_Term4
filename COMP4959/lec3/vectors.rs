struct VecIter {
    v: Vec<i32>,
    i: usize,
}

impl VecIter{
  fn new(v: Vec<i32>) -> Self {
    Self{v, i:0}
  }
}

impl Iterator for VecIter {
    type Item = i32;

    fn next(&mut self) -> Option<i32> {
      if self.i >= self.v.len() {
        None
      } else {
        let res = self.v[self.i];
        self.i += 1;
        Some(res)
      }
    }
}

fn main() {
  let v = vec![1, 2, 3, 4, 5]; 
  let mut it = VecIter::new(v);
  while let Some(x) = it.next() {
    println!("{x}");
  }
}