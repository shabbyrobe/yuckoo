package yuckoo

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
)

var source = rng(0)

func genRand(i uint64, max uint64) uint64 {
	return source.Uint64() % max
}

func genSeq(i uint64, max uint64) uint64 {
	return i
}

func TestCuckooAddHas(t *testing.T) {
	const resultHit, resultFull, resultMiss = 1, 2, 3

	tw := tabwriter.NewWriter(os.Stdout, 4, 2, 2, ' ', 0)
	defer tw.Flush()

	for idx, tc := range []struct {
		capacity int
		items    uint64
		gen      func(i, max uint64) uint64
	}{
		{1000, 100, genSeq},
		{100, 1000, genSeq},
		{1000, 1000, genSeq},
		{1000, 1001, genSeq},
		{1000, 2000, genSeq},

		{1000, 100, genRand},
		{100, 1000, genRand},
		{1000, 1000, genRand},
		{1000, 1001, genRand},
		{1000, 2000, genRand},
		{10000, 20000, genRand},
		{100000, 200000, genRand},
	} {

		t.Run(fmt.Sprintf("%d-%d,%d", idx, tc.items, tc.capacity), func(t *testing.T) {
			f := New(tc.capacity)
			values := make([]uint64, tc.items)
			for i := uint64(0); i < tc.items; i++ {
				v := tc.gen(i, tc.items)
				f.Add(v)
				values[i] = v
			}

			var results = map[uint64]int{}
			hit, miss, full := 0, 0, 0
			for _, v := range values {
				if has, isFull := f.Has(v); has {
					results[v] = resultHit
					hit++
				} else if isFull {
					results[v] = resultFull
					full++
				} else {
					results[v] = resultMiss
					miss++
				}
			}

			min := uint64(tc.capacity)
			if tc.items < min {
				min = tc.items
			}
			hitRate := float64(hit) / float64(min)

			enc := f.Encode(nil)

			fmt.Fprintf(tw,
				"n: %d"+"\t"+
					"cap: %d"+"\t"+
					"hit: %d"+"\t"+
					"miss: %d"+"\t"+
					"full: %d"+"\t"+
					"hit: %.02f"+"\t"+
					"sz: %d"+"\n",
				tc.items, tc.capacity, hit, miss, full, hitRate, len(enc))

			if miss > 0 {
				t.Fatal("miss")
			}
			if hitRate < 0.99 {
				t.Fatal("hitrate", hitRate)
			}

			dec, left, err := Decode(enc)
			if err != nil {
				t.Fatal(err)
			}
			if len(left) > 0 {
				t.Fatal()
			}

			for _, v := range values {
				r := 0
				if has, isFull := dec.Has(v); has {
					r = resultHit
				} else if isFull {
					r = resultFull
				} else {
					r = resultMiss
				}
				if r != results[v] {
					t.Fatal()
				}
			}
		})
	}
}
