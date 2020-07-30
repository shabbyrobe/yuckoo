package yuckoo

var (
	masks   = [256]uint64{} // only require 65 entires, but 256 prevents a bounds check
	altHash = [256]uint64{}
)

func init() {
	for i := 0; i <= 64; i++ {
		masks[i] = (1 << i) - 1
	}
	for i := uint64(0); i < 256; i++ {
		altHash[i] = fnv1aU64(i)
	}
}
