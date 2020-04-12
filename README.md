Yuckoo - Cuckoo Filter For My Specific Use Case
===============================================

Package cuckoo contains a specialised version of a cuckoo filter that does not
support deletion and uses fnv1-a instead of MetroHash as ARM is a target
platform. MetroHash does not have a Go implementation that outperforms fnv1-a
on ARM.

This may be updated to use
[FxHash](https://docs.rs/hashers/1.0.0/hashers/fx_hash/struct.FxHasher64.html)
at some point, which appears to outperform fnv1-a for most use-cases.


## Expectation Management

I recommend using https://github.com/seiflotfy/cuckoofilter instead of this,
unless you are absolutely sure you don't need deletion and you are absolutely
sure you aren't worried by the following caveats.

Feel free to use this or take bits from it as you see fit (MIT license == go
nuts). I recommend vendoring it and testing it thoroughly if you intend to use
it as I will change this without warning to suit my own needs at any time.

Issues may be responded to whenever I happen to get around to them, but PRs are
unlikely to be accepted.
