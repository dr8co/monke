package lexer

import "testing"

// BenchmarkLexer measures the performance of the lexer by tokenizing a sample program
func BenchmarkLexer(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(input)
		for {
			tok := l.NextToken()
			if tok.Type == "EOF" {
				break
			}
		}
	}
}

// BenchmarkLexerLarge measures the performance of the lexer with a larger input
func BenchmarkLexerLarge(b *testing.B) {
	// Create a larger input by repeating the sample program
	input := `
let five = 5;
let ten = 10;
let add = fn(x, y) {
    x + y;
};
let result = add(five, ten);
`
	largeInput := ""
	for i := 0; i < 100; i++ {
		largeInput += input
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(largeInput)
		for {
			tok := l.NextToken()
			if tok.Type == "EOF" {
				break
			}
		}
	}
}
