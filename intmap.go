package intmap

/*
Package intmap provides a persistent integer map, based on Chris Okasaki's
Fast Mergeable Integer Maps (http://ittc.ku.edu/~andygill/papers/IntMap98.pdf).
The datastructure is based on a big-endian Patricia tree.

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

type leaf struct {
	k uint64
	v interface{}
}
