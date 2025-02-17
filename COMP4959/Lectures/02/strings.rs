fn main() {
  let s1 = "hello"; // type &str
  // println!("{}", s1[0]);
  println!("{}", &s1[0..1]);
  let s2 = String::from("hello");
  println!("{} {} {}", s, s.len(), s.capacity());
  print_str(&s2);
  print_string_ref(&s2);
  print_string(s2);
}

fn print_str(s: &str) {
  println!("{}", s);
}

fn print_string(s: String) {
  println!("{}", s);
}

fn print_string_ref(s: &String) {
  println!("{}", s);
}
