fn main() {
    println!("Simple Calculator in Rust");

    // Test the calculator functions
    let a = 10;
    let b = 5;
    let c = 0;

    println!("{} + {} = {}", a, b, add(a, b));
    println!("{} - {} = {}", a, b, subtract(a, b));
    println!("{} * {} = {}", a, b, multiply(a, b));
    println!("{} / {} = {}", a, b, divide(a, b));
    println!("{} % {} = {}", a, b, modulo(a, b));

    // Test error handling
    println!("\nTesting error cases:");
    println!("{} / {} = {}", a, c, divide(a, c));
    println!("{} % {} = {}", a, c, modulo(a, c));
}

// Addition function
fn add(a: i32, b: i32) -> i32 {
    a + b
}

// Subtraction function
fn subtract(a: i32, b: i32) -> i32 {
    a - b
}

// Multiplication function
fn multiply(a: i32, b: i32) -> i32 {
    a * b
}

// Division function
fn divide(a: i32, b: i32) -> i32 {
    if b == 0 {
        println!("Error: Division by zero!");
        0
    } else {
        a / b
    }
}

// Modulo function
fn modulo(a: i32, b: i32) -> i32 {
    if b == 0 {
        println!("Error: Division by zero in modulo operation!");
        0
    } else {
        a % b
    }
}