package intmap

/*
Package intmap provides a persistent integer map, based on Chris Okasaki's
Fast Mergeable Integer Maps (http://ittc.ku.edu/~andygill/papers/IntMap98.pdf).
The datastructure is a big-endian Patricia tree.

In addition to having lookup and insert times that are at least competitive with
a typical binary tree, an IntMap boasts exceptionally fast merge performance.
*/

// IntMaps can hold any type of value.
type Value interface{}

// Internal interface; abstracts differences between branches and leaves. This
// effectively trades time (extra work to deref interface calls) for space (an
// alternative implementation would probably have to include space for both
// branch and leaf data). Shame that Go doesn't have tagged unions.
type node interface {
	Lookup(k uint64) (v Value, exists bool)
}

type branch struct {
	p uint64
	m uint64
	t [2]node
}

type emptyBranch struct{}

var empty = emptyBranch{}

type leaf struct {
	k uint64
	v interface{}
}

func (b *branch) Get(k uint64) (Value, bool) {
	// Simple comparison is possible because big-endian Patricia trees are teh awesome
	if k <= b.p {
		if b.t[0] == nil {
			return nil, false
		} else {
			return b.t[0].Lookup(k)
		}
	} else {
		if b.t[1] == nil {
			return nil, false
		} else {
			return b.t[1].Lookup(k)
		}
	}
}

func (b *branch) join(p0, p1 uint64, t0, t1 node) *branch {
	var (
		m = branchingBit(p0, p1)
		t [2]node
	)

	if zeroBit(p0, m) {
		t[0], t[1] = t0, t1
	} else {
		t[0], t[1] = t1, t0
	}

	return &branch{
		p: mask(p0, m),
		m: m,
		t: [2]node{t0, t1},
	}
}

func branchingBit(p0, p1 uint64) uint64 {
	return (p0 ^ p1) & (^(p0 ^ p1))
}

func zeroBit(k, m uint64) bool {
	return k&m == 0
}

func mask(k, m uint64) uint64 {
	return k & (m - 1)
}

func matchPrefix(k, p, m uint64) bool {
	return mask(k, m) == p
}

func (l leaf) Lookup(k uint64) (Value, bool) {
	if k == l.k {
		return l.v, true
	} else {
		return nil, false
	}
}
