// Command profile runs a Monkey program and measures its execution time.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/dr8co/monke/evaluator"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/object"
	"github.com/dr8co/monke/parser"
)

var (
	cpuprofile   = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile   = flag.String("memprofile", "", "write memory profile to file")
	traceprofile = flag.String("trace", "", "write execution trace to file")
	program      = flag.String("program", "fibonacci", "program to profile (fibonacci, factorial, array, hash, complex)")
)

// Sample Monkey programs for profiling
var programs = map[string]string{
	"fibonacci": `
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
		fibonacci(20);
	`,
	"factorial": `
		let factorial = fn(n) {
			if (n == 0) {
				return 1;
			} else {
				return n * factorial(n - 1);
			}
		};
		factorial(10);
	`,
	"array": `
		let map = fn(arr, f) {
			let iter = fn(arr, accumulated) {
				if (len(arr) == 0) {
					return accumulated;
				} else {
					return iter(rest(arr), push(accumulated, f(first(arr))));
				}
			};
			return iter(arr, []);
		};

		let reduce = fn(arr, initial, f) {
			let iter = fn(arr, result) {
				if (len(arr) == 0) {
					return result;
				} else {
					return iter(rest(arr), f(result, first(arr)));
				}
			};
			return iter(arr, initial);
		};

		let arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
		let double = fn(x) { return x * 2; };
		let sum = fn(x, y) { return x + y; };

		let doubled = map(arr, double);
		reduce(doubled, 0, sum);
	`,
	"hash": `
		let hash = {"one": 1, "two": 2, "three": 3, "four": 4, "five": 5};
		let sum = 0;

		let keys = ["one", "two", "three", "four", "five"];
		let i = 0;

		while (i < len(keys)) {
			let key = keys[i];
			sum = sum + hash[key];
			i = i + 1;
		}

		sum;
	`,
	"complex": `
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

		let map = fn(arr, f) {
			let iter = fn(arr, accumulated) {
				if (len(arr) == 0) {
					return accumulated;
				} else {
					return iter(rest(arr), push(accumulated, f(first(arr))));
				}
			};
			return iter(arr, []);
		};

		let reduce = fn(arr, initial, f) {
			let iter = fn(arr, result) {
				if (len(arr) == 0) {
					return result;
				} else {
					return iter(rest(arr), f(result, first(arr)));
				}
			};
			return iter(arr, initial);
		};

		let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
		let compute = fn(x) { return fibonacci(x); };
		let sum = fn(x, y) { return x + y; };

		let results = map(numbers, compute);
		reduce(results, 0, sum);
	`,
}

// builtins is a map of built-in functions that are available to the Monkey program
var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return &object.Null{}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `first` not supported, got %s", args[0].Type())}
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					return arg.Elements[length-1]
				}
				return &object.Null{}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `last` not supported, got %s", args[0].Type())}
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1)
					copy(newElements, arg.Elements[1:length])
					return &object.Array{Elements: newElements}
				}
				return &object.Null{}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `rest` not supported, got %s", args[0].Type())}
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newElements := make([]object.Object, length+1)
				copy(newElements, arg.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `push` not supported, got %s", args[0].Type())}
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return &object.Null{}
		},
	},
}

func main() {
	flag.Parse()
	exit := func(code int) {
		os.Exit(code)
	}

	// Set up CPU profiling if requested
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
			if err != nil {
				return
			}
			exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)
		if err := pprof.StartCPUProfile(f); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "could not start CPU profile: %v\n", err)
			if err != nil {
				return
			}
			exit(2)
		}
		defer pprof.StopCPUProfile()
	}

	// Set up execution tracing if requested
	if *traceprofile != "" {
		f, err := os.Create(*traceprofile)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "could not create trace profile: %v\n", err)
			if err != nil {
				return
			}
			exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)
		if err := trace.Start(f); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "could not start trace profile: %v\n", err)
			if err != nil {
				return
			}
			exit(1)
		}
		defer trace.Stop()
	}

	// Get the program to profile
	input, ok := programs[*program]
	if !ok {
		_, err := fmt.Fprintf(os.Stderr, "unknown program: %s\n", *program)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(os.Stderr, "available programs: fibonacci, factorial, array, hash, complex\n")
		if err != nil {
			return
		}
		exit(1)
	}

	// Create an environment with built-in functions
	env := object.NewEnvironment()
	for name, builtin := range builtins {
		env.Set(name, builtin)
	}

	// Run the program and measure the time
	start := time.Now()

	// Lexing
	l := lexer.New(input)

	// Parsing
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		_, err := fmt.Fprintf(os.Stderr, "parser errors:\n")
		if err != nil {
			return
		}
		for _, msg := range p.Errors() {
			_, err := fmt.Fprintf(os.Stderr, "\t%s\n", msg)
			if err != nil {
				return
			}
		}
		exit(1)
	}

	// Evaluation
	result := evaluator.Eval(program, env)

	elapsed := time.Since(start)
	fmt.Printf("Program: %s\n", *program)
	fmt.Printf("Result: %s\n", result.Inspect())
	fmt.Printf("Time: %s\n", elapsed)

	// Write the memory profile if requested
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "could not create memory profile: %v\n", err)
			if err != nil {
				return
			}
			exit(1)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "could not write memory profile: %v\n", err)
			if err != nil {
				return
			}
			exit(1)
		}
	}
}
