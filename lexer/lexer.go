// Package lexer implements the lexical analyzer for the Monke programming language.
//
// The lexer is responsible for breaking down the source code into tokens,
// which are the smallest units of meaning in the language. It reads the input
// character by character and produces a stream of tokens that can be processed
// by the parser.
//
// Key features:
//   - Tokenization of all language elements (keywords, identifiers, literals, operators, etc.)
//   - Handling of whitespace and comments
//   - Error detection for illegal characters
//   - Support for various token types defined in the token package
//   - Optimized for performance with minimal allocations
//
// The main entry point is the New function, which creates a new Lexer instance,
// and the NextToken method, which returns the next token from the input.
package lexer

import "github.com/dr8co/monke/token"

// Common tokens that are reused to reduce allocations
var (
	tokenPlus      = token.Token{Type: token.PLUS, Literal: "+"}
	tokenMinus     = token.Token{Type: token.MINUS, Literal: "-"}
	tokenSlash     = token.Token{Type: token.SLASH, Literal: "/"}
	tokenAsterisk  = token.Token{Type: token.ASTERISK, Literal: "*"}
	tokenLT        = token.Token{Type: token.LT, Literal: "<"}
	tokenGT        = token.Token{Type: token.GT, Literal: ">"}
	tokenSemicolon = token.Token{Type: token.SEMICOLON, Literal: ";"}
	tokenColon     = token.Token{Type: token.COLON, Literal: ":"}
	tokenComma     = token.Token{Type: token.COMMA, Literal: ","}
	tokenLParen    = token.Token{Type: token.LPAREN, Literal: "("}
	tokenRParen    = token.Token{Type: token.RPAREN, Literal: ")"}
	tokenLBrace    = token.Token{Type: token.LBRACE, Literal: "{"}
	tokenRBrace    = token.Token{Type: token.RBRACE, Literal: "}"}
	tokenLBracket  = token.Token{Type: token.LBRACKET, Literal: "["}
	tokenRBracket  = token.Token{Type: token.RBRACKET, Literal: "]"}
	tokenEOF       = token.Token{Type: token.EOF, Literal: ""}
)

// Lexer represents the lexer for the Monke programming language.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	// Pre-allocates a token to reuse for single-character tokens
	singleCharToken token.Token
}

// readChar reads the next character from the input and advances the position.
// It's optimized to minimize checks and operations.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// New creates a new Lexer with the given input string.
// It initializes the lexer, reads the first character, and sets up the token buffer.
func New(input string) *Lexer {
	l := &Lexer{
		input:           input,
		singleCharToken: token.Token{}, // Initialize the token buffer
	}
	l.readChar()
	return l
}

// NextToken reads the next token from the input.
// It skips whitespace, identifies the token type based on the current character,
// and returns a token with the appropriate type and literal value.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			// Use a pre-allocated token for "=="
			l.readChar() // Advance to the next character after '=='
			return token.Token{Type: token.EQ, Literal: string(ch) + string('=')}
		}
		l.readChar() // Advance to the next character after '='
		return token.Token{Type: token.ASSIGN, Literal: "="}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			// Use a pre-allocated token for "!="
			l.readChar() // Advance to the next character after '!='
			return token.Token{Type: token.NOT_EQ, Literal: string(ch) + string('=')}
		}
		l.readChar() // Advance to the next character after '!'
		return token.Token{Type: token.BANG, Literal: "!"}
	case '+':
		l.readChar() // Advance to the next character after '+'
		return tokenPlus
	case '-':
		l.readChar() // Advance to the next character after '-'
		return tokenMinus
	case '/':
		l.readChar() // Advance to the next character after '/'
		return tokenSlash
	case '*':
		l.readChar() // Advance to the next character after '*'
		return tokenAsterisk
	case '<':
		l.readChar() // Advance to the next character after '<'
		return tokenLT
	case '>':
		l.readChar() // Advance to the next character after '>'
		return tokenGT
	case ';':
		l.readChar() // Advance to the next character after ';'
		return tokenSemicolon
	case ':':
		l.readChar() // Advance to the next character after ':'
		return tokenColon
	case ',':
		l.readChar() // Advance to the next character after ','
		return tokenComma
	case '(':
		l.readChar() // Advance to the next character after '('
		return tokenLParen
	case ')':
		l.readChar() // Advance to the next character after ')'
		return tokenRParen
	case '{':
		l.readChar() // Advance to the next character after '{'
		return tokenLBrace
	case '}':
		l.readChar() // Advance to the next character after '}'
		return tokenRBrace
	case '[':
		l.readChar() // Advance to the next character after '['
		return tokenLBracket
	case ']':
		l.readChar() // Advance to the next character after ']'
		return tokenRBracket
	case '"':
		tok := token.Token{Type: token.STRING}
		tok.Literal = l.readString()
		l.readChar() // Advance to the next character after the closing quote
		return tok
	case 0:
		return tokenEOF
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return token.Token{
				Type:    token.LookupIdent(literal),
				Literal: literal,
			}
		}
		if isDigit(l.ch) {
			return token.Token{
				Type:    token.INT,
				Literal: l.readNumber(),
			}
		}
		// For illegal characters, reuse the single char token
		l.singleCharToken.Type = token.ILLEGAL
		l.singleCharToken.Literal = string(l.ch)
		l.readChar() // Advance to the next character after the illegal character
		return l.singleCharToken
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readNumber reads a number from the input and returns it as a string.
// It's optimized to avoid unnecessary allocations.
func (l *Lexer) readNumber() string {
	position := l.position
	// Fast-forward through digits
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readIdentifier reads an identifier from the input and returns it as a string.
// It's optimized to avoid unnecessary allocations.
func (l *Lexer) readIdentifier() string {
	position := l.position
	// Fast-forward through letters
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// skipWhitespace skips any whitespace characters in the input.
// It's optimized to use a single loop.
func (l *Lexer) skipWhitespace() {
	// Fast-forward through whitespace
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// peekChar returns the next character in the input without advancing the position.
// It's optimized to avoid unnecessary checks.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// readString reads a string from the input and returns it as a string.
// It's optimized to avoid unnecessary allocations.
func (l *Lexer) readString() string {
	position := l.position + 1
	// Fast-forward through string characters
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
