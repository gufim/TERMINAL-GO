package main

import "testing"

func BenchmarkNewTappableSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// TODO: wywołaj newTappableSelect(...)
		_ = i
	}
}

func BenchmarkNewCustomEntry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// TODO: wywołaj NewCustomEntry(...)
		_ = i
	}
}
