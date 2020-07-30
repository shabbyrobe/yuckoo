package yuckoo

const offset64 = 14695981039346656037
const prime64 = 1099511628211

func fnv1aU64(v uint64) uint64 {
	return (offset64 ^ v) * prime64
}

func fnv1a64(b []byte) uint64 {
	var hash uint64 = offset64
	for _, c := range b {
		hash ^= uint64(c)
		hash *= prime64
	}
	return hash
}

func getNextPow2(n uint64) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return int(n)
}

/*
This is a fixed-increment version of Java 8's SplittableRandom generator
See http://dx.doi.org/10.1145/2714064.2660195 and
http://docs.oracle.com/javase/8/docs/api/java/util/SplittableRandom.html

It is a very fast generator passing BigCrush, and it can be useful if
for some reason you absolutely want 64 bits of state; otherwise, we
rather suggest to use a xoroshiro128+ (for moderately parallel
computations) or xorshift1024* (for massively parallel computations)
generator.

http://xoshiro.di.unimi.it/splitmix64.c
*/
type rng uint64

func (sm64 *rng) Int63() int64 { return int64(sm64.Uint64() >> 1) }

func (sm64 *rng) Uint64() uint64 {
	*sm64 += 0x9E3779B97F4A7C15
	var z uint64 = uint64(*sm64)
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}
