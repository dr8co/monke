# Monke Profiling Tool

This tool is used to profile the Monke interpreter with various sample programs to identify performance bottlenecks.

## Usage

```bash
# Build the profiling tool
cd cmd/profile
go build

# Run with default settings (fibonacci program)
./profile

# Run with CPU profiling
./profile -cpuprofile=cpu.prof

# Run with memory profiling
./profile -memprofile=mem.prof

# Run with execution tracing
./profile -trace=trace.out

# Run a different program
./profile -program=factorial
./profile -program=array
./profile -program=hash
./profile -program=complex
```

## Available Programs

- `fibonacci`: Calculates the 20th Fibonacci number using recursion
- `factorial`: Calculates the factorial of 10 using recursion
- `array`: Demonstrates array operations with map and reduce functions
- `hash`: Demonstrates hash table operations
- `complex`: A complex program that combines multiple features

## Analyzing Profiling Results

To analyze CPU profiling results:
```bash
go tool pprof -top cpu.prof
```

To analyze memory profiling results:
```bash
go tool pprof -top mem.prof
```

To view the execution trace:
```bash
go tool trace trace.out
```

## Results

The profiling results are documented in the `/profile_results.md` file at the root of the project.