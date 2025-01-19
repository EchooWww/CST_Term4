use std::io::{self};
use std::fmt;

struct Student {
  first_name: String,
  last_name: String,
  score: i32
}
impl Student {
  fn new (first_name: &str, last_name: &str, score: i32) -> Student{
    Student {first_name: first_name.to_string(), last_name:last_name.to_string(), score}  }
}
type StudentList = Vec<Student>;


impl fmt::Display for Student {
  fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
    write!(f, "- {}, {}", self.last_name, self.first_name)
  }}

fn read_students() -> io::Result<StudentList> {
  let mut students = Vec::new();
  for line in io::stdin().lines() {
    let line = line?;
    let info:Vec<&str> = line.trim().split(" ").collect();
    let first_name = info[0];
    let last_name = info[1];
    let age = info[2].parse().unwrap();
    students.push(Student::new(first_name, last_name, age));
  }
  Ok(students)
}

fn print_summary(students: &[Student]){
  let num = students.len();
  if num == 0 {
    println!("No student record found.");
    return;
  }

  let sum: i32 = students.iter().map(|s| s.score).sum();
  let min = students.iter().map(|s| s.score).min().unwrap();
  let max = students.iter().map(|s| s.score).max().unwrap();
  let average = sum as f32 / num as f32;

  println!("\nnumber of students: {}", num);
  println!("average score: {:.1}", average);
  if min == max {
    println!("minimum score = maximum score = {}", min);
    for student in students.iter().filter(|s| s.score == min) {
        println!("{student}");
    }
  } else {
      println!("minimum score = {}", min);
      for student in students.iter().filter(|s| s.score == min) {
          println!("{student}");
      }
      println!("maximum score = {}", max);
      for student in students.iter().filter(|s| s.score == max) {
          println!("{student}");
      }
  }
}


fn main() {
  let students = read_students().unwrap();
  print_summary(&students);
}
