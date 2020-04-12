package cuckoo

import (
	"math/bits"
)

const bucketSize = 4
const maxKicks = 500

type fingerprint uint8

type bucket [bucketSize]fingerprint

type Filter struct {
	buckets   []bucket
	bucketPow uint8
	full      bool
	rng       rng
}

func New(capacity int) *Filter {
	capacity = getNextPow2(uint64(capacity)) / bucketSize
	if capacity <= 0 {
		panic("next power of 2 after capacity must yield a number > 0")
	}

	buckets := make([]bucket, capacity)
	return &Filter{
		buckets,
		uint8(bits.TrailingZeros(uint(capacity))),
		false,
		rng(offset64),
	}
}

func (f *Filter) Add(v uint64) bool {
	var bucket *bucket

	h1 := fnv1aU64(v)
	fp := fingerprint(h1%255 + 1)
	idx1 := h1 & masks[f.bucketPow]
	bucket = &f.buckets[idx1]

	// NOTE: we do not insert duplicate fingerprints because we do not need
	// to support deletion:
	if bucket[0] == 0 || bucket[0] == fp {
		bucket[0] = fp
		return true
	}
	if bucket[1] == 0 || bucket[1] == fp {
		bucket[1] = fp
		return true
	}
	if bucket[2] == 0 || bucket[2] == fp {
		bucket[2] = fp
		return true
	}
	if bucket[3] == 0 || bucket[3] == fp {
		bucket[3] = fp
		return true
	}

	h2 := h1 ^ fnv1aU64(uint64(fp))
	idx2 := h2 & masks[f.bucketPow]
	bucket = &f.buckets[idx2]
	if bucket[0] == 0 || bucket[0] == fp {
		bucket[0] = fp
		return true
	}
	if bucket[1] == 0 || bucket[1] == fp {
		bucket[1] = fp
		return true
	}
	if bucket[2] == 0 || bucket[2] == fp {
		bucket[2] = fp
		return true
	}
	if bucket[3] == 0 || bucket[3] == fp {
		bucket[3] = fp
		return true
	}

	bucketIdx := idx1
	if f.rng.Uint64()&0b1 == 1 {
		bucketIdx = idx2
	}

	var randn uint64
	var randl int
	var n int
	for n = 0; n < maxKicks; n++ {
		if randl == 0 {
			randn, randl = f.rng.Uint64(), 64
		}
		fpIdx := randn & 0b11
		randn >>= 2
		randl -= 2

		prev := fp
		fp = f.buckets[bucketIdx][fpIdx]
		f.buckets[bucketIdx][fpIdx] = prev

		mask := masks[f.bucketPow]
		bucketIdx = (bucketIdx ^ altHash[fp]) & mask

		bucket = &f.buckets[bucketIdx]
		if bucket[0] == 0 || bucket[0] == fp {
			bucket[0] = fp
			return true
		}
		if bucket[1] == 0 || bucket[1] == fp {
			bucket[1] = fp
			return true
		}
		if bucket[2] == 0 || bucket[2] == fp {
			bucket[2] = fp
			return true
		}
		if bucket[3] == 0 || bucket[3] == fp {
			bucket[3] = fp
			return true
		}
	}

	f.full = true
	return false
}

func (f *Filter) Has(v uint64) (ok bool, full bool) {
	h1 := fnv1aU64(v)
	fp := fingerprint(h1%255 + 1) // '0' fingerprint means 'not set'
	idx1 := h1 & masks[f.bucketPow]
	b := f.buckets[idx1]
	if b[0] == fp || b[1] == fp || b[2] == fp || b[3] == fp {
		return true, false
	}

	h2 := h1 ^ fnv1aU64(uint64(fp))
	idx2 := h2 & masks[f.bucketPow]
	b = f.buckets[idx2]
	if b[0] == fp || b[1] == fp || b[2] == fp || b[3] == fp {
		return true, false
	}

	return false, f.full
}
