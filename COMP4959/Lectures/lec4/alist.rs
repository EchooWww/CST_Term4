struct Alist<K, V>(Vec<(K,V)>);

// we need <K, V> after impl to indicate they are type variables instead of an actual type
impl <K, V> Alist<K, V>  {
  fn new() -> Self {
    Self(Vec::new())
  }

  fn len(&self) -> usize {
    self.0.len()
  }

  fn add(&mut self, k: K, v: V) {
    self.0.push((k, v));
  }

  fn into_iter(self) -> std::vec::IntoIter<(K, V)> {
    self.0.into_iter()
  }

  fn iter(&self) -> std::slice::Iter<(K, V)> {
    self.0.iter()
  }
}

impl<K, V> IntoIterator for Alist<K, V> {
  type Item = (K, V);
  type IntoIter = std::vec::IntoIter<(K, V)>;
  fn into_iter(self)->Self::IntoIter {
    Alist::into_iter(self)
  }
}

// to be able to do for x in &alist (iter() returns a reference to the alist)
// lifetime parameter: the item must live as long as the alist

impl <'a, K,V> IntoIterator for &'a Alist<K, V> {
    type Item = &'a (K, V);
    type IntoIter = std::slice::Iter<'a, (K, V)>;
    fn into_iter(self) -> Self::IntoIter {
      self.iter()
    }
}


// struct AlistIter<K, V> {
//   l: Alist<K, V>,
//   i: usize,
// }

// impl <K:Clone, V:Clone> Iterator for AlistIter<K, V> {
//   type Item = (K, V);

//   fn next(&mut self) -> Option<(K, V)> {
//     if self.i >= self.l.len() {
//       None
//     } else {
//       let res = self.l.0[self.i].clone();
//       self.i += 1;
//       Some(res)
//     }
//   }
// }

fn main() {
  let mut l = Alist::new();
  l.add(1, "hello");
  l.add(2, "world");
  l.add(3, "foo");

  // l.into_iter().for_each(|(k, v)| {
  //  println!("{}: {}", k, v);
  //});
  l.iter().for_each(|(k, v)| {
    println!("{}: {}", k, v);
  });
  for (k, v) in &l {
    println!("{}: {}", k, v);
  }
}
