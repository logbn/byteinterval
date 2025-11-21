package byteinterval

import (
	"bytes"
	"sync"

	"github.com/rdleal/intervalst/interval"
)

// New returns a new byte interval search tree.
func New[V any]() *Tree[V] {
	return &Tree[V]{
		tree: interval.NewMultiValueSearchTreeWithOptions[*Interval[V], []byte](
			bytes.Compare,
			interval.TreeWithIntervalPoint(),
		),
	}
}

// Tree is a thread safe interval radix tree.
type Tree[V any] struct {
	tree *interval.MultiValueSearchTree[*Interval[V], []byte]
	mut  sync.RWMutex
}

// Insert makes the value available between interval [start, end). Exact matches on end are excluded.
func (t *Tree[V]) Insert(start, end []byte, val V) (i *Interval[V]) {
	if len(end) == 0 {
		end = start
	}
	if bytes.Compare(start, end) > 0 {
		return
	}
	i = &Interval[V]{t, start, end, val, false}
	t.mut.Lock()
	defer t.mut.Unlock()
	t.tree.Insert(start, end, i)
	return
}

// Find returns value of all intervals intersecting point key.
func (t *Tree[V]) Find(k []byte) (vals []V) {
	t.mut.RLock()
	defer t.mut.RUnlock()
	intervals, _ := t.tree.AllIntersections(k, k)
	m := map[*Interval[V]]bool{}
	for _, i := range intervals {
		if bytes.Equal(k, i.end) && !bytes.Equal(k, i.start) {
			continue
		}
		if _, ok := m[i]; !ok {
			vals = append(vals, i.val)
			m[i] = true
		}
	}
	return
}

// FindAny returns value of all intervals intersecting any point in keys (no dupes).
func (t *Tree[V]) FindAny(keys ...[]byte) (vals []V) {
	t.mut.RLock()
	defer t.mut.RUnlock()
	m := map[*Interval[V]]bool{}
	for _, k := range keys {
		intervals, _ := t.tree.AllIntersections(k, k)
		for _, i := range intervals {
			if bytes.Equal(k, i.end) && !bytes.Equal(k, i.start) {
				continue
			}
			if _, ok := m[i]; !ok {
				vals = append(vals, i.val)
				m[i] = true
			}
		}
	}
	return
}

func (t *Tree[V]) remove(i *Interval[V]) {
	t.mut.Lock()
	defer t.mut.Unlock()
	found, ok := t.tree.Find(i.start, i.end)
	if !ok {
		return
	}
	for n, i2 := range found {
		if i == i2 {
			found = append(found[:n], found[n+1:]...)
			break
		}
	}
	if len(found) == 0 {
		t.tree.Delete(i.start, i.end)
	} else {
		t.tree.Upsert(i.start, i.end, found...)
	}
	i.rem = true
}

// Interval represnts an interval inserted into the tree.
type Interval[V any] struct {
	tree  *Tree[V]
	start []byte
	end   []byte
	val   V
	rem   bool
}

// Remove is the only way to remove items from the tree.
func (i *Interval[V]) Remove() {
	if i.rem {
		return
	}
	i.tree.remove(i)
}
