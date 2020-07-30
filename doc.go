// Package yuckoo contains a specialised version of a cuckoo filter that
// does not support deletion and uses fnv1-a instead of MetroHash as
// ARM is a target platform. MetroHash does not have a Go implementation that
// outperforms fnv1-a on ARM.
package yuckoo
