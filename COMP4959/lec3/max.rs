fn max(a: &i32, b:&i32) -> &i32{
    if a > b {
        a
    } else {
        b
    }
}

fn main(){
    let a = 10;
    let b = 20;
    let m = max(&a, &b);
    println!("{}", m);
}