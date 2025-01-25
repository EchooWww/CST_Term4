use std::env;
use std::io::{self};

enum Range {
  Single(usize),
  From(usize),
  Between(usize, usize),
  To(usize)
}

impl Range {
  fn parse(s: &str) -> Option<Range> {
    if let Some((start,end)) = s.split_once('-'){
      if start.is_empty() {
        // -M
        return end.parse().ok().map(Range::To);
      } 
      if end.is_empty() {
        // N-
        return start.parse().ok().map(Range::From);
      } 
      let start = start.parse().ok()?;
      let end = end.parse().ok()?;
      if start > end {
        return None;
          // N-M
      }
      return Some(Range::Between(start, end));
    }
    s.parse::<usize>().ok().map(Range::Single)
  }

  fn contains(&self, n: usize) -> bool {
    match self {
        Range::Single(x) => *x == n,
        Range::Between(start, end) => n >= *start && n <= *end,
        Range::From(start) => n >= *start,
        Range::To(end) => n <= *end,
    }
  }
}

fn main() {
  let args: Vec<String> = env::args().collect();
  if args.len() != 2 {
    println!("Wrong command. Usage: ./cut -c<ranges> or -f<ranges>");
    return;
  }
  
  let param = &args[1];
  let mut ranges: Vec<Range> = Vec::new();

  if !(param.starts_with("-c") || param.starts_with("-f")) {
    println!("Invalid argument: must start with -c or -f");
    return;
  }    
  if param[2..].is_empty() {
    println!("Error: No ranges specified.");
    return; 
  }
  for r in param[2..].split(',') {
    match Range::parse(r) {
        Some(range) => ranges.push(range),
        None => {
          eprintln!("Error: Invalid range"); 
          return; 
        }
    }
};

for line in io::stdin().lines() {
    let line = line.unwrap();
    if param.starts_with("-c") {
        // Character mode
        let result: String = line
            .chars()
            .enumerate()
            .filter(|(i, _)| ranges.iter().any(|r| r.contains(i + 1)))
            .map(|(_, c)| c)
            .collect();
        println!("{}", result);
    } else {
        // Field mode (tab-separated fields)
        let result: Vec<&str> = line
            .split('\t')
            .enumerate()
            .filter(|(i, _)| ranges.iter().any(|r| r.contains(i + 1)))
            .map(|(_, field)| field)
            .collect();
        println!("{}", result.join("\t"));
    }
  }
}