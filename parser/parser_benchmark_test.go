package parser

import (
	"testing"

	"github.com/dr8co/monke/lexer"
)

// benchmarkParser is a helper function for benchmarking the parser
func benchmarkParser(input string, b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := lexer.New(input)
		p := New(l)
		p.ParseProgram()
	}
}

// BenchmarkParseSimpleExpression measures the performance of parsing a simple expression
func BenchmarkParseSimpleExpression(b *testing.B) {
	input := `
	let x = 5;
	let y = 10;
	let z = x + y;
	`
	benchmarkParser(input, b)
}

// BenchmarkParseComplexExpression measures the performance of parsing a complex expression
func BenchmarkParseComplexExpression(b *testing.B) {
	input := `
	let x = 5;
	let y = 10;
	let z = 8;
	let result = (x + y) * z / (x + y - z) + x * y;
	`
	benchmarkParser(input, b)
}

// BenchmarkParseFunction measures the performance of parsing a function definition
func BenchmarkParseFunction(b *testing.B) {
	input := `
	let add = fn(x, y) {
		return x + y;
	};
	`
	benchmarkParser(input, b)
}

// BenchmarkParseRecursiveFunction measures the performance of parsing a recursive function
func BenchmarkParseRecursiveFunction(b *testing.B) {
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
	`
	benchmarkParser(input, b)
}

// BenchmarkParseArrayAndHash measures the performance of parsing array and hash literals
func BenchmarkParseArrayAndHash(b *testing.B) {
	input := `
	let myArray = [1, 2, 3, 4, 5];
	let myHash = {"one": 1, "two": 2, "three": 3};
	`
	benchmarkParser(input, b)
}

// BenchmarkParseLargeProgram measures the performance of parsing a larger program
func BenchmarkParseLargeProgram(b *testing.B) {
	input := `
	let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x + y;
	};
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 == 10;
	10 != 9;

	"foobar"
	"foo bar"
	[1, 2];
	{"foo": "bar"}

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
	`
	benchmarkParser(input, b)
}
