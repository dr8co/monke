package evaluator

import (
	"testing"

	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/object"
	"github.com/dr8co/monke/parser"
)

// benchmarkEval is a helper function for benchmarking the evaluator
func benchmarkEval(input string, b *testing.B) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eval(program, env)
	}
}

// BenchmarkIntegerArithmetic measures the performance of integer arithmetic operations
func BenchmarkIntegerArithmetic(b *testing.B) {
	input := `
	let x = 5;
	let y = 10;
	x + y;
	x - y;
	x * y;
	y / x;
	-x;
	`
	benchmarkEval(input, b)
}

// BenchmarkStringConcatenation measures the performance of string concatenation
func BenchmarkStringConcatenation(b *testing.B) {
	input := `
	let x = "Hello";
	let y = " World";
	x + y;
	x + x + y + y;
	`
	benchmarkEval(input, b)
}

// BenchmarkArrayLiteral measures the performance of array creation and access
func BenchmarkArrayLiteral(b *testing.B) {
	input := `
	let x = [1, 2, 3, 4, 5];
	let y = x[2];
	let z = x[4];
	x[0] + x[1] + y + x[3] + z;
	`
	benchmarkEval(input, b)
}

// BenchmarkHashLiteral measures the performance of hash creation and access
func BenchmarkHashLiteral(b *testing.B) {
	input := `
	let x = {"one": 1, "two": 2, "three": 3, "four": 4, "five": 5};
	let y = x["three"];
	let z = x["five"];
	x["one"] + x["two"] + y + x["four"] + z;
	`
	benchmarkEval(input, b)
}

// BenchmarkFunctionCall measures the performance of function calls
func BenchmarkFunctionCall(b *testing.B) {
	input := `
	let adder = fn(a, b) { a + b };
	let multiplier = fn(a, b) { a * b };
	adder(5, 10) + multiplier(2, 4);
	`
	benchmarkEval(input, b)
}

// BenchmarkRecursiveFunction measures the performance of recursive function calls
func BenchmarkRecursiveFunction(b *testing.B) {
	input := `
	let fibonacci = fn(x) {
		if (x == 0) {
			return 0;
		} else {
			if (x == 1) {
				return 1;
			} else {
				return fibonacci(x - 1) + fibonacci(x - 2);
			}
		}
	};
	fibonacci(10);
	`
	benchmarkEval(input, b)
}

// BenchmarkConditionals measures the performance of conditional expressions
func BenchmarkConditionals(b *testing.B) {
	input := `
	let x = 10;
	let y = 5;
	if (x > y) {
		x;
	} else {
		y;
	}
	if (x < y) {
		x;
	} else {
		y;
	}
	if (x == y) {
		x;
	} else {
		y;
	}
	`
	benchmarkEval(input, b)
}

// BenchmarkComplexExpression measures the performance of a complex expression
func BenchmarkComplexExpression(b *testing.B) {
	input := `
	let x = 5;
	let y = 10;
	let z = 8;
	let result = (x + y) * z / (x + y - z) + x * y;
	result;
	`
	benchmarkEval(input, b)
}
