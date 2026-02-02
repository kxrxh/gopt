package goptions

import (
	"testing"
)

func BenchmarkSome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Some(42)
	}
}

func BenchmarkNone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = None[int]()
	}
}

func BenchmarkMap(b *testing.B) {
	o := Some(21)
	fn := func(x int) int { return x * 2 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Map(o, fn)
	}
}

func BenchmarkAndThen(b *testing.B) {
	o := Some(10)
	fn := func(x int) Option[int] { return Some(x + 1) }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AndThen(o, fn)
	}
}

func BenchmarkUnwrapOr(b *testing.B) {
	o := Some(42)
	defaultVal := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = o.UnwrapOr(defaultVal)
	}
}

func BenchmarkUnwrapOrNone(b *testing.B) {
	o := None[int]()
	defaultVal := 42
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = o.UnwrapOr(defaultVal)
	}
}

func BenchmarkFilter(b *testing.B) {
	o := Some(4)
	pred := func(x int) bool { return x%2 == 0 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = o.Filter(pred)
	}
}

func BenchmarkFlatten(b *testing.B) {
	o := Some(Some(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Flatten(o)
	}
}

func BenchmarkMatch(b *testing.B) {
	o := Some(42)
	onSome := func(x int) int { return x }
	onNone := func() int { return 0 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Match(o, onSome, onNone)
	}
}
