# Monke Programming Language Specification

This document describes the syntax and semantics of the Monke programming language.

## 1. Introduction

Monke is a simple, interpreted programming language with a C-like syntax.

It is dynamically typed and supports first-class functions, closures,
and basic data structures like arrays and hash maps.

## 2. Lexical Elements

### 2.1 Comments

Monke does not currently support comments.

### 2.2 Identifiers

Identifiers start with a letter or underscore and can contain letters, digits, and underscores.

```
identifier = letter { letter | digit | "_" } .
letter = "a"..."z" | "A"..."Z" | "_" .
digit = "0"..."9" .
```

### 2.3 Keywords

The following keywords are reserved and cannot be used as identifiers:

```
fn    let    true    false    if    else    return
```

### 2.4 Operators and Delimiters

The following characters and character sequences represent operators and delimiters:

```
+    -    *    /    =    ==    !=    <    >    !
(    )    {    }    [    ]    ,    ;    :
```

### 2.5 Literals

#### 2.5.1 Integer Literals

Integer literals consist of a sequence of digits.

```
integer = digit { digit } .
```

#### 2.5.2 String Literals

String literals are enclosed in double quotes.

```
string = '"' { character } '"' .
```

#### 2.5.3 Boolean Literals

Boolean literals are `true` and `false`.

#### 2.5.4 Array Literals

Array literals are enclosed in square brackets and contain a comma-separated list of expressions.

```
array = "[" [ expression { "," expression } ] "]" .
```

#### 2.5.5 Hash Literals

Hash literals are enclosed in curly braces and contain a comma-separated list of key-value pairs.

```
hash = "{" [ expression ":" expression { "," expression ":" expression } ] "}" .
```

## 3. Types

Monke has the following built-in types:

- Integer: 64-bit signed integer
- Boolean: true or false
- String: sequence of characters
- Array: ordered collection of values
- Hash: collection of key-value pairs
- Function: first-class function
- Null: represents the absence of a value

## 4. Expressions

### 4.1 Primary Expressions

#### 4.1.1 Identifiers

Identifiers refer to variables, functions, or built-in functions.

#### 4.1.2 Literals

Literals represent fixed values.

#### 4.1.3 Parenthesized Expressions

Expressions can be enclosed in parentheses to control precedence.

```
( expression )
```

### 4.2 Function Literals

Function literals define anonymous functions.

```
fn ( parameters ) { statements }
```

### 4.3 Call Expressions

Call expressions invoke functions.

```
expression ( arguments )
```

### 4.4 Index Expressions

Index expressions access elements of arrays or hashes.

```
expression [ expression ]
```

### 4.5 Prefix Expressions

Prefix expressions apply an operator to a single operand.

```
operator expression
```

Supported prefix operators:
- `-`: Negation (for integers)
- `!`: Logical NOT (for booleans)

### 4.6 Infix Expressions

Infix expressions apply an operator to two operands.

```
expression operator expression
```

Supported infix operators:
- `+`: Addition (for integers and strings)
- `-`: Subtraction (for integers)
- `*`: Multiplication (for integers)
- `/`: Division (for integers)
- `<`: Less than (for integers)
- `>`: Greater than (for integers)
- `==`: Equal to (for all types)
- `!=`: Not equal to (for all types)

### 4.7 If Expressions

If expressions provide conditional evaluation.

```
if ( expression ) { statements } [ else { statements } ]
```

## 5. Statements

### 5.1 Expression Statements

Expression statements evaluate an expression and discard the result.

```
expression ;
```

### 5.2 Let Statements

Let statements bind a value to an identifier.

```
let identifier = expression ;
```

### 5.3 Return Statements

Return statements return a value from a function.

```
return expression ;
```

### 5.4 Block Statements

Block statements group multiple statements together.

```
{ statements }
```

## 6. Built-in Functions

Monke provides the following built-in functions:

- `len(arg)`: Returns the length of a string or array
- `first(array)`: Returns the first element of an array
- `last(array)`: Returns the last element of an array
- `rest(array)`: Returns a new array containing all elements except the first
- `push(array, element)`: Returns a new array with the element added to the end
- `puts(args...)`: Prints the arguments to the console

## 7. Evaluation Rules

Monke uses eager evaluation.

Expressions are evaluated from left to right,
with operator precedence determining the order of operations.

## 8. Scoping Rules

Monke uses lexical scoping.

Variables are visible within the block where they are defined and any nested blocks,
unless shadowed by a variable with the same name in a nested block.

## 9. Error Handling

Monke does not have explicit error handling mechanisms like try/catch.
Runtime errors result in error objects that terminate execution.