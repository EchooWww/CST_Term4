use std::collections::HashMap;
fn primes (m:usize, n:usize) -> Vec<usize>{
  if n < 2 {
    return vec![];
  }
  let mut is_prime = vec![true; n+1];
  is_prime[0] = false;
  is_prime[1] = false;
  for i in 2..=((n as f64).sqrt() as usize) {
    if is_prime[i] {
      for p in (i*i..=n).step_by(i) {
        is_prime[p] = false;
      }
    }
  }
  is_prime
    .iter()
    .enumerate()
    .filter(|&(i, &is_p)|is_p&&i>=m)
    .map(|(i, _)|i)
    .collect()
}

fn sort_digits(num:&usize)->String {
  let mut digits: Vec<char> = num.to_string().chars().collect();
  digits.sort();
  digits.into_iter().collect()
}

fn most_permutations(nums: &[usize]) -> usize {
  let mut groups: HashMap<String, Vec<&usize>> = HashMap::new();
  for num in nums {
      let key = sort_digits(num); 
      groups.entry(key).or_insert_with(Vec::new).push(num);
  }

  groups.values().map(|group| group.len()).max().unwrap_or(0)
}


fn main() {
  let v = primes(100000, 999999);
  println!("Total number of 6-digit primes: {:?}", v.len());
  let most = most_permutations(&v);
  println!("Size of largest 6-digit primes permutation set: {most}");
}