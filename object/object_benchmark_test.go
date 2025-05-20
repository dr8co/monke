package object

import (
	"testing"
)

// BenchmarkStringHashKey measures the performance of the string hash key calculation
func BenchmarkStringHashKey(b *testing.B) {
	s := &String{Value: "Hello World"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.HashKey()
	}
}

// BenchmarkIntegerHashKey measures the performance of the integer hash key calculation
func BenchmarkIntegerHashKey(b *testing.B) {
	integer := &Integer{Value: 42}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = integer.HashKey()
	}
}

// BenchmarkBooleanHashKey measures the performance of the boolean hash key calculation
func BenchmarkBooleanHashKey(b *testing.B) {
	b1 := &Boolean{Value: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = b1.HashKey()
	}
}

// BenchmarkHashCreation measures the performance of creating a hash with multiple entries
func BenchmarkHashCreation(b *testing.B) {
	keys := []*String{
		{Value: "one"},
		{Value: "two"},
		{Value: "three"},
		{Value: "four"},
		{Value: "five"},
	}
	values := []*Integer{
		{Value: 1},
		{Value: 2},
		{Value: 3},
		{Value: 4},
		{Value: 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pairs := make(map[HashKey]HashPair)
		for j := 0; j < 5; j++ {
			key := keys[j]
			value := values[j]
			hashed := key.HashKey()
			pairs[hashed] = HashPair{Key: key, Value: value}
		}
		_ = &Hash{Pairs: pairs}
	}
}

// BenchmarkHashLookup measures the performance of looking up values in a hash
func BenchmarkHashLookup(b *testing.B) {
	keys := []*String{
		{Value: "one"},
		{Value: "two"},
		{Value: "three"},
		{Value: "four"},
		{Value: "five"},
	}
	values := []*Integer{
		{Value: 1},
		{Value: 2},
		{Value: 3},
		{Value: 4},
		{Value: 5},
	}

	pairs := make(map[HashKey]HashPair)
	for j := 0; j < 5; j++ {
		key := keys[j]
		value := values[j]
		hashed := key.HashKey()
		pairs[hashed] = HashPair{Key: key, Value: value}
	}
	hash := &Hash{Pairs: pairs}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 5; j++ {
			key := keys[j]
			hashed := key.HashKey()
			_ = hash.Pairs[hashed]
		}
	}
}
