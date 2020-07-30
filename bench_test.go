package yuckoo

import (
	"crypto/rand"
	"io"
	"testing"
)

var BenchHasResult, BenchFullResult bool

func BenchmarkAdd(b *testing.B) {
	const capacity = 10000
	filter := New(capacity)

	j := uint64(0)
	end := uint64(capacity * 2)

	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		filter.AddUint64(j)
		j++
		if j > end {
			j = 0
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	const cap = 10000
	filter := New(cap)

	for i := uint64(0); i < uint64(10000); i++ {
		filter.AddUint64(i)
	}

	b.ResetTimer()
	j := uint64(0)
	for i := 0; i < b.N; i++ {
		BenchHasResult, BenchFullResult = filter.HasUint64(j)
		j++
		if j > 20000 {
			j = 0
		}
	}
}

func BenchmarkAgainstOtherLib_Insert(b *testing.B) {
	// Tries to get a benchmark result that is comparable to the work
	// the benchmarks are doing in here: https://github.com/seiflotfy/cuckoofilter

	const cap = 10000
	filter := New(cap)

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.AddBytes(hash[:])
	}
}

func BenchmarkAgainstOtherLib_Lookup(b *testing.B) {
	const cap = 10000
	filter := New(cap)

	var hash [32]byte
	for i := 0; i < 10000; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.AddBytes(hash[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.HasBytes(hash[:])
	}
}
