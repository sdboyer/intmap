package intmap

/*
Package intmap provides a persistent integer map, based on Chris Okasaki's
Fast Mergeable Integer Maps (http://ittc.ku.edu/~andygill/papers/IntMap98.pdf).

In addition to having lookup and insert times that are at least competitive with
a typical binary tree, an IntMap boasts exceptionally fast merge performance.
*/

type IntMap struct {
}

type branch struct {
	prefix    uint8
	branchbit uint8
}

type leaf struct {
	key uint64
	val interface{}
}
