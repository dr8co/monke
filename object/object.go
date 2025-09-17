// Package object defines the object system for the Monke programming language.
//
// This package implements the runtime object system that represents values
// during the execution of a Monke program. It defines various types of objects
// such as integers, booleans, strings, arrays, hashes, functions, and built-ins.
//
// Key components:
//   - Object interface: The base interface for all runtime values
//   - Various object types (Integer, Boolean, String, Array, Hash, Function, etc.)
//   - Environment: Stores variable bindings during execution
//   - Hashable interface: For objects that can be used as hash keys
//   - Optimized hash table implementation with key caching for better performance
//
// The evaluator uses the object system to represent and manipulate values
// during program execution.
package object

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/dr8co/monke/ast"
)

//nolint:revive
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

// ObjectType represents the type of object.
type ObjectType string

// Object is the interface that wraps the basic operations of all Monke objects.
// All Monke objects implement this interface.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents a Monke integer value.
type Integer struct {
	Value int64
}

// Type returns the type of the object.
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Inspect returns a string representation of the object.
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }

// Boolean represents a Monke boolean value.
type Boolean struct {
	Value bool
}

// Type returns the type of the object.
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Inspect returns a string representation of the object.
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }

// String represents a Monke string value.
type String struct {
	Value string
	// Cache for the hash key to avoid recalculating it
	hashKey *HashKey
}

// Type returns the type of the object.
func (s *String) Type() ObjectType { return STRING_OBJ }

// Inspect returns a string representation of the object.
func (s *String) Inspect() string { return s.Value }

// Null represents a Monke null value.
type Null struct{}

// Type returns the type of the object.
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Inspect returns a string representation of the object.
func (n *Null) Inspect() string { return "null" }

// ReturnValue represents a Monke return value.
type ReturnValue struct {
	Value Object
}

// Type returns the type of the object.
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Inspect returns a string representation of the object.
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error represents a Monke error.
type Error struct {
	Message string
}

// Type returns the type of the object.
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Inspect returns a string representation of the object.
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function represents a Monke function.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns the type of the object.
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// Inspect returns a string representation of the object.
func (f *Function) Inspect() string {
	var out strings.Builder
	params := make([]string, 0, len(f.Parameters))

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// BuiltinFunction represents a Monke builtin function.
type BuiltinFunction func(args ...Object) Object

// Builtin represents a Monke builtin.
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns the type of the object.
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Inspect returns a string representation of the object.
func (b *Builtin) Inspect() string { return "builtin function" }

// Array represents a Monke array.
type Array struct {
	Elements []Object
}

// Type returns the type of the object.
func (a *Array) Type() ObjectType { return ARRAY_OBJ }

// Inspect returns a string representation of the object.
func (a *Array) Inspect() string {
	var out strings.Builder

	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e.Inspect()
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey represents a hash key.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashKey returns the hash key for the object.
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns the hash key for the object.
func (i *Integer) HashKey() HashKey {
	//nolint:gosec
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey returns the hash key for the object.
func (s *String) HashKey() HashKey {
	// Return the cached hash key if available
	if s.hashKey != nil {
		return *s.hashKey
	}

	// Calculate the hash key
	h := fnv.New64a()
	_, err := h.Write([]byte(s.Value))
	if err != nil {
		return HashKey{Type: ERROR_OBJ, Value: 0}
	}

	// Create and cache the hash key
	hashKey := HashKey{Type: s.Type(), Value: h.Sum64()}
	s.hashKey = &hashKey
	return hashKey
}

// HashPair represents a hash pair.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents a Monke hash.
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the type of the object.
func (h *Hash) Type() ObjectType { return HASH_OBJ }

// Inspect returns a string representation of the object.
func (h *Hash) Inspect() string {
	var out strings.Builder

	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Hashable represents an object that can be used as a hash key.
type Hashable interface {
	HashKey() HashKey
}
