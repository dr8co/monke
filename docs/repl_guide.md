# Monke REPL User Guide

This guide explains how to use the Monke Read-Eval-Print Loop (REPL) interface,
which allows you to interactively write and execute Monke code.

## Getting Started

### Starting the REPL

To start the Monke REPL, run the `monke` executable:

```bash
monke
```

You'll see a welcome message and a prompt (`>>`) where you can start entering Monke code.

## REPL Interface

The REPL interface consists of the following elements:

1. **Title Bar**: Displays "Monkey Programming Language REPL" at the top of the interface.
2. **Welcome Message**: Shows a personalized greeting with your username.
3. **Command History**: Displays previously entered commands and their results.
4. **Input Prompt**: The `>>` prompt where you enter Monke code.
5. **Help Text**: Shows available keyboard shortcuts at the bottom of the interface.

## Basic Usage

### Entering Code

Type your Monke code at the prompt and press Enter to execute it:
  
```console
>> let x = 5;
```

The REPL will evaluate your code and display the result:

```console
>> let x = 5;
5
```

### Viewing Results

Results are displayed in green text below your input. If there's an error, it will be displayed in red text.

### Multi-line Input

You can enter multi-line expressions. The REPL will evaluate the code when you complete the expression:

```console
>> if (x > 3) {
     x + 10
   } else {
     x - 10
   }
15
```

### Command History

The REPL keeps track of your command history during the current session.
Previously entered commands and their results remain visible above the current prompt.

## Features

### Variables and Expressions

You can define variables and use them in expressions:

```console
>> let name = "World";
>> let greeting = "Hello, " + name + "!";
>> greeting
"Hello, World!"
```

### Functions

Define and call functions:

```console
>> let add = fn(a, b) { a + b };
>> add(5, 10)
15
```

### Data Structures

Create and manipulate arrays and hash maps:

```console
>> let myArray = [1, 2, 3, 4, 5];
>> myArray[2]
3

>> let myHash = {"name": "Monke", "type": "Language"};
>> myHash["name"]
"Monke"
```

### Built-in Functions

Use Monke's built-in functions:

```console
>> len("hello")
5

>> puts("Hello, World!")
Hello, World!
null
```

## Keyboard Shortcuts

- **Enter**: Execute the current input
- **Esc** or **Ctrl+C**: Exit the REPL

## Tips and Tricks

1. **Persistent Environment**: Variables and functions defined in the REPL persist for the duration of the session.

2. **Semicolons**: Statements should end with a semicolon (`;`).

3. **Expression Results**: The REPL displays the result of the last expression evaluated.

4. **Error Messages**: If your code has syntax errors, the REPL will display detailed error messages to help you fix the issues.

5. **Experimenting**: The REPL is perfect for experimenting with language features and testing small code snippets before incorporating them into larger programs.

## Example Session

Here's an example of a typical REPL session:

```console
>> let x = 10;
10

>> let y = 5;
5

>> x + y
15

>> let max = fn(a, b) { if (a > b) { a } else { b } };
fn(a, b) {
if (a > b) { a } else { b }
}

>> max(x, y)
10

>> let array = [1, 2, 3, 4, 5];
[1, 2, 3, 4, 5]

>> len(array)
5

>> array[2]
3
```

## Conclusion

The Monke REPL is a powerful tool for learning the language, testing code snippets, and experimenting with different features. It provides immediate feedback, making it ideal for both beginners and experienced users.

For more information about the Monke language itself, refer to the language specification document.