// Package token defines the token types and structures for the Monke programming language.
//
// Tokens are the smallest units of meaning in the language, produced by the lexer
// during the lexical analysis phase. Each token represents a specific language element
// such as a keyword, identifier, operator, or delimiter.
//
// Key components:
//   - TokenType: A type representing different categories of tokens
//   - Token: A structure containing the type and literal value of a token
//   - Constants for all token types supported by the language
//   - Lookup functions for identifying keywords
//
// This package is used primarily by the lexer to categorize input text and by the
// parser to understand the structure of the program.
package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers & literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	COLON     = ":"
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks if the given identifier is a keyword.
// If it is, it returns the corresponding token type.
// Otherwise, it returns the IDENT token type.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
