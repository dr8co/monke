// Example 2: Conditionals and Boolean Operations

// Define some variables
let x = 10;
let y = 5;
let z = 10;

// Boolean literals
let isTrue = true;
let isFalse = false;

// Comparison operators
puts("Comparison Operators:");
puts("x == y:", x == y);  // false
puts("x != y:", x != y);  // true
puts("x > y:", x > y);    // true
puts("x < y:", x < y);    // false
puts("x == z:", x == z);  // true

// Boolean operators
puts("\nBoolean Operators:");
puts("!isTrue:", !isTrue);   // false
puts("!isFalse:", !isFalse); // true

// If expressions
puts("\nIf Expressions:");

let result = if (x > y) { "x is greater than y" } else { "x is not greater than y" };
puts(result); // "x is greater than y"

// Nested if expressions
let grade = 85;
let letterGrade = if (grade >= 90) {
    "A"
} else if (grade >= 80) {
    "B"
} else if (grade >= 70) {
    "C"
} else if (grade >= 60) {
    "D"
} else {
    "F"
};

puts("Grade:", grade);
puts("Letter Grade:", letterGrade); // "B"

// If expressions without else
let message = if (x == z) { "x equals z" };
puts(message); // "x equals z"

// If expressions can be used in other expressions
let value = 5 + if (isTrue) { 10 } else { 0 };
puts("Value:", value); // 15