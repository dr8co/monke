# Monke

Monke is a simple, interpreted programming language implemented in Go.

This project is inspired by the book "Writing An Interpreter In Go" by Thorsten Ball.

It features a REPL (Read-Eval-Print Loop), a lexer, parser, AST, and an evaluator, all written from scratch.

## Features

- **Lexer**: Tokenizes input source code.
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens.
- **AST**: Represents the structure of parsed code.
- **Evaluator**: Executes the AST, supporting variables, functions, and basic data types.
- **REPL**: Interactive shell for running Monke code.
- **Built-in Functions**: Includes basic built-in functions for convenience.

## Getting Started

### Prerequisites

- Go 1.24 or newer

### Installation

To install Monke, you can either clone the repository and build it yourself or install it directly using Go.

#### Option 1: Clone and Build

1. Clone the repository:

   ```sh
   git clone https://github.com/dr8co/monke.git
   cd monke
   go build -o monke
   ```

2. Run the binary:

   ```sh
    ./monke
    ```

#### Option 2: Install with Go

```sh
go install github.com/dr8co/monke@latest
```

You should see a prompt where you can enter Monke code interactively.

## Project Structure

- `main.go` — Entry point, starts the REPL.
- `cmd/profile/` — Profiling tool for performance analysis.
- `lexer/` — Lexical analyzer (tokenizer).
- `parser/` — Parser for Monke language.
- `ast/` — Abstract Syntax Tree definitions.
- `object/` — Object system and environment.
- `evaluator/` — Evaluates the AST.
- `repl/` — REPL implementation.
- `token/` — Token definitions.
- `docs/` — Documentation and tasks.

## Example Usage

```console
>> let x = 5;
>> x + 10;
15
>> fn add(a, b) { a + b; }
>> add(2, 3);
5
```

## Testing

To run the tests:

```sh
go test ./...
```

## Profiling

Monke includes a profiling tool to analyze performance:

```sh
# Build and run the profiling tool
cd cmd/profile
go build
./profile
```

For more details on profiling options and available test programs, see the [profiling tool README](cmd/profile/README.md).

## Contributing

Contributions are welcome! Please open issues or submit pull requests for improvements or bug fixes.

## License

This project is licensed under the MIT License.
