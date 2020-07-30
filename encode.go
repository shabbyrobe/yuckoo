package yuckoo

import (
	"fmt"
)

const version = 1

func (f *Filter) Encode(into []byte) (b []byte) {
	into = into[:cap(into)]

	// 4 bytes for LE length prefix, 1 byte for version, 1 byte for full flag + pow, 2
	// bytes reserved
	off := 8

	sz := off + (len(f.buckets) * bucketSize)
	if len(into) < sz {
		into = make([]byte, sz)
	}

	full := uint8(0)
	if f.full {
		full = 1
	}

	ln := sz - 4
	_ = into[8]
	into[0] = byte(ln)
	into[1] = byte(ln >> 8)
	into[2] = byte(ln >> 16)
	into[3] = byte(ln >> 24)
	into[4] = version
	into[5] = full<<7 | f.bucketPow
	into[6] = 0
	into[7] = 0

	idx := off
	for _, bucket := range f.buckets {
		into[idx+0] = byte(bucket[0])
		into[idx+1] = byte(bucket[1])
		into[idx+2] = byte(bucket[2])
		into[idx+3] = byte(bucket[3])
		idx += 4
	}

	return into[:sz]
}

func Decode(raw []byte) (filt *Filter, left []byte, err error) {
	if len(raw) < 6 {
		return nil, raw, fmt.Errorf("cuckoo: incomplete filter")
	}

	_ = raw[8]

	// size does not include the length of its own value:
	sz := uint32(raw[0]) | uint32(raw[1])<<8 | uint32(raw[2])<<16 | uint32(raw[3])<<24
	if uint32(len(raw)) < sz+4 {
		return nil, raw, fmt.Errorf("cuckoo: incomplete filter")
	}

	if raw[4] != version {
		return nil, raw, fmt.Errorf("cuckoo: unknown version")
	}

	full := raw[5]&0b_1000_0000 != 0
	pow := raw[5] & 0b_0111_1111

	// account for 4 bytes of metadata after the length:
	capacity := (sz - 4) / bucketSize
	if (sz-4)%bucketSize != 0 {
		return nil, raw, fmt.Errorf("cuckoo: invalid length")
	}

	buckets := make([]bucket, capacity)
	x := 8
	for i := uint32(0); i < capacity; i++ {
		buckets[i][0] = fingerprint(raw[x+0])
		buckets[i][1] = fingerprint(raw[x+1])
		buckets[i][2] = fingerprint(raw[x+2])
		buckets[i][3] = fingerprint(raw[x+3])
		x += 4
	}

	flt := &Filter{
		buckets:   buckets,
		bucketPow: pow,
		full:      full,
		rng:       rng(offset64),
	}
	return flt, raw[x:], nil
}
