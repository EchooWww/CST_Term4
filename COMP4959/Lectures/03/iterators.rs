fn main() {
  let v = vec![1,2,3];
  v.into_iter().for_each(|x| println!("{}", x));
  // println!("{:?}", v); // v has moved
  v.iter().for_each(|x| println!("{}", x));
  println!("{:?}", v); // v has not moved
  let mut = vec![1,2,3];
  v.iter_mut().for_each(|x| *x += 1);

}