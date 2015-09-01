package intmap

/*
Package intmap provides a persistent integer map, based on Chris Okasaki's
Fast Mergeable Integer Maps (http://ittc.ku.edu/~andygill/papers/IntMap98.pdf).
The datastructure is a big-endian Patricia tree.

In addition to having lookup and insert times that are at least competitive with
a typical binary tree, an IntMap boasts exceptionally fast merge performance.
*/

// IntMaps can hold any type of value.
type Entry interface{}

type Node struct {
	// The prefix
	p uint64
	// The branching bit
	m uint64
	// flag indicating whether the Node is non-empty. negation so that zero-value is correct
	notEmpty bool
	// The children
	t [2]*Node
	// The value (iff a leaf)
	l *leaf
}

type leaf struct {
	k uint64
	v Entry
}

var root = &Node{}

// NewMap creates a new empty IntMap for use. This is operation performs no allocations.
func NewMap() *Node {
	return root
}

func (n *Node) Insert(k uint64, v Entry) *Node {
	var newNode = &Node{
		p:        n.p,
		m:        n.m,
		notEmpty: true,
		l:        &leaf{k, v},
	}
	if !n.notEmpty {
		return newNode
	} else {
		if n.l != nil {
			// On a leaf, so replace if key is same, else join
			if n.l.k == k {
				return newNode
			} else {
				n.join(k, n.l.k, newNode, n)
			}
		} else {
			// On a branch, so we must recurse
			if mask(k, n.m) == n.p {
				// There's a prefix match, so recurse down
			}
		}
	}
}

func (b *Node) Get(k uint64) (Entry, bool) {
	// Simple comparison is possible because big-endian Patricia trees are teh awesome
	if k <= b.p {
		if b.t[0] == nil {
			return nil, false
		} else {
			return b.t[0].Get(k)
		}
	} else {
		if b.t[1] == nil {
			return nil, false
		} else {
			return b.t[1].Get(k)
		}
	}
}

func (b *Node) join(p0, p1 uint64, t0, t1 *Node) *Node {
	var (
		m = branchingBit(p0, p1)
		t [2]*Node
	)

	if zeroBit(p0, m) {
		t[0], t[1] = t0, t1
	} else {
		t[0], t[1] = t1, t0
	}

	return &Node{
		p: mask(p0, m),
		m: m,
		t: t,
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
