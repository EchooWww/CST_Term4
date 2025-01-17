use std::io::{self, BufRead, BufReader};
use std::fs::File;

// // Keep reading lines from stdin until EOF
// fn main() {
//   let stdin = io::stdin();
//   let mut line = String::new();
//   // read_lines stores the line terminator
//   while let Ok(n) = stdin.read_line(&mut line) {
//     if n ==0 {
//       break;
//     } 
//     print!("{line}");
//     line.clear();
//   }
// }

fn read_lines() -> io::Result<Vec<String>> {
  let mut line = String::new();
  let mut v = Vec::new();
  while io::stdin().read_line(&mut line)? > 0 {
    v.push(line.clone());
    line.clear();
  }
  Ok(v)
}

fn read_file_lines(file: &str) -> io::Result<Vec<String>> {
  let f = File::open(file)?;
  let mut reader = BufReader::new(f);
  let mut line = String::new();
  let mut v = Vec::new();
  while reader.read_line(&mut line)? > 0 {
    v.push(line.clone());
    line.clear();
  }
  Ok(v)
}

fn lines () -> io::Result<Vec<String>> {
  io::stdin().lines().collect()
}

fn file_lines(file: &str) -> io::Result<Vec<String>> {
  BufReader::new(File::open(file)?).lines().collect()
}

fn main() {
  // let lines = read_lines();
  // println!("{:?}", lines);
  // let file = "read_lines.rs";
  // let lines = read_file_lines(file);
  // println!("{:?}", lines);
  let lines: Vec<String> = file_lines("read_lines.rs");
  println!("{:?}", lines);
}