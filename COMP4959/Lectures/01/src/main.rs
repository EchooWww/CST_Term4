mod numbers; 
mod arrays;
mod vectors;
mod slices;

mod tuples {
    pub fn run() {
        let t = (1,2);
        let (x,y) = t;
        println!("{:?}", t); //this will throw an error because the ownership of t has been moved to x and y
    }
}

fn main() {
    numbers::run();
    arrays::run();
    vectors::run();
    slices::run();
}

fn double(a: &mut [i32]) {
    for x in a { // x has type &mut i32
        // the for loop is creating an iterator, which is immutable, so we need to dereference x to make it mutable
        *x *= 2;
    }
}