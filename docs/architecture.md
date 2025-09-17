# Monke Interpreter: Architecture and Design Decisions

This document describes the architecture of the Monke programming language
interpreter and explains the key design decisions made during its implementation.

## Overall Architecture

The Monke interpreter follows a classic interpreter design pattern with the following pipeline:

```txt
Source Code → Lexer → Parser → AST → Evaluator → Result
```

Each component has a specific responsibility:

1. **Lexer**: Converts source code text into tokens
2. **Parser**: Transforms tokens into an Abstract Syntax Tree (AST)
3. **Evaluator**: Executes the AST to produce results

The interpreter is implemented as a tree-walking interpreter,
which means it directly traverses and evaluates the AST without any intermediate representation.

## Component Breakdown

### Lexer (`lexer` package)

The lexer is responsible for breaking down the source code into tokens.
It reads the input character by character and produces a stream of tokens.

**Design Decisions:**

- **Simple Character-by-Character Scanning**: The lexer uses a straightforward approach of reading one character at a time, which makes the code easy to understand and maintain.
- **Look-Ahead**: The lexer uses a one-character look-ahead to handle multi-character tokens like `==` and `!=`.
- **No Regex**: The lexer avoids using regular expressions, making it more portable and easier to understand.
- **Token Types**: Tokens are categorized into different types (keywords, identifiers, literals, operators, etc.) to simplify parsing.

### Parser (`parser` package)

The parser converts the token stream into an Abstract Syntax Tree (AST).
It implements a recursive descent parser with Pratt parsing (precedence climbing) for expressions.

**Design Decisions:**

- **Pratt Parsing**: The parser uses Pratt parsing (also known as precedence climbing) for handling expressions with different operator precedences. This approach is both powerful and relatively simple to implement.
- **Recursive Descent**: For statements and other language constructs, the parser uses recursive descent, which closely mirrors the grammar of the language.
- **Error Reporting**: The parser collects errors during parsing rather than stopping at the first error, allowing it to report multiple issues at once.
- **Prefix and Infix Functions**: The parser uses maps of prefix and infix parsing functions to handle different types of expressions, making it easy to extend with new expression types.

### AST (`ast` package)

The Abstract Syntax Tree (AST) represents the structure of the program.
Each node in the tree corresponds to a language construct (expression, statement, etc.).

**Design Decisions:**

- **Node Interface**: All AST nodes implement a common `Node` interface, allowing them to be treated uniformly.
- **Expression and Statement Interfaces**: Nodes are further categorized as either expressions (which produce values) or statements (which perform actions), reflecting the language's grammar.
- **String Representation**: Each node can produce a string representation of itself, which is useful for debugging and testing.
- **Immutable Nodes**: AST nodes are designed to be immutable, simplifying the evaluation process.

### Evaluator (`evaluator` package)

The evaluator traverses the AST and executes the program. It implements the semantics of the Monke language.

**Design Decisions:**

- **Tree-Walking Interpreter**: The evaluator directly traverses and evaluates the AST without any intermediate representation, prioritizing simplicity over performance.
- **Environment-Based Scoping**: Variable scopes are implemented using environment objects that can be nested to support lexical scoping.
- **First-Class Functions**: Functions are treated as first-class values, allowing them to be passed around, returned from other functions, and stored in variables.
- **Closures**: Functions capture their defining environment, enabling closures.
- **Error Handling**: Errors are represented as values that can be passed around, allowing for consistent error handling throughout the evaluation process.
- **Built-in Functions**: Common functions are provided as built-ins, implemented directly in Go rather than in Monke.

### Object System (`object` package)

The object system defines the runtime values that can exist in a Monke program.
It includes integers, booleans, strings, arrays, hashes, functions, and more.

**Design Decisions:**

- **Object Interface**: All runtime values implement a common `Object` interface, allowing them to be treated uniformly.
- **Value Representation**: Each type of value has its own struct with appropriate fields.
- **Hashable Interface**: Objects that can be used as hash keys implement the `Hashable` interface, allowing them to be used in hash maps.
- **Environment**: The environment is implemented as a map from strings (variable names) to objects, with support for nested environments.

### REPL (`repl` package)

The REPL (Read-Eval-Print Loop) provides an interactive interface for users to enter Monke code and see the results immediately.

**Design Decisions:**

- **Modern Terminal UI**: The REPL uses the Charm libraries (Bubbletea, Bubbles, and Lipgloss) to create a modern, user-friendly terminal interface.
- **Command History**: The REPL keeps track of command history, allowing users to see their previous inputs and results.
- **Styled Output**: Different types of output (results, errors) are styled differently for better readability.
- **Persistent Environment**: The environment persists across commands, allowing users to define variables and functions that can be used in later commands.

## Key Design Principles

### Simplicity Over Performance

The Monke interpreter prioritizes simplicity and clarity of implementation over raw performance.
This makes it an excellent learning tool and a good foundation for experimentation.

### Modularity

The interpreter is designed with a clear separation of concerns,
with each component having a well-defined responsibility.

This makes the code easier to understand, test, and modify.

### Extensibility

The design allows for easy extension with new language features.
For example, adding a new expression type requires changes to the parser and evaluator,
but the overall architecture remains the same.

### Error Handling

Errors are treated as values that can be passed around, rather than using exceptions.
This leads to more explicit error handling and makes the control flow easier to follow.

## Future Directions

While the current implementation is a tree-walking interpreter, future versions could include:

1. **Bytecode Compiler**: Compiling the AST to bytecode and executing it on a virtual machine would significantly improve performance.
2. **Just-In-Time Compilation**: JIT compilation could provide even better performance for hot code paths.
3. **Static Type Checking**: Adding optional static type checking would catch more errors at compile time.
4. **Module System**: A proper module system would allow for better code organization and reuse.
5. **Standard Library**: Expanding the built-in functions into a comprehensive standard library would make the language more useful for real-world tasks.

## Conclusion

The Monke interpreter demonstrates a clean, modular design that prioritizes simplicity and clarity.
Its architecture follows established patterns for interpreter implementation,
making it both a good learning tool and a solid foundation for experimentation with language design and implementation.
