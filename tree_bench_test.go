package byteinterval

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	b.Run("Overlap Small", func(b *testing.B) {
		benchInsert(b, testGenKeys)
	})
	b.Run("Overlap Big", func(b *testing.B) {
		benchInsert(b, testGenKeysFull)
	})
}

func benchInsert(b *testing.B, fn func(int) [][][]byte) {
	for _, n := range []int{1 << 8, 1 << 12, 1 << 16} {
		keys := fn(n)
		tree := New[*int]()
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for b.Loop() {
				for j, k := range keys {
					tree.Insert(k[0], k[1], &j)
				}
			}
		})
	}
}

func testGenKeys(n int) [][][]byte {
	res := make([][][]byte, n)
	for i := 0; i < n; i++ {
		start := rand.Int63()
		res[i] = [][]byte{
			[]byte(fmt.Sprintf("%016x", start)),
			[]byte(fmt.Sprintf("%016x", start+10000)),
		}
	}
	return res
}

func testGenKeysFull(n int) [][][]byte {
	res := make([][][]byte, n)
	for i := 0; i < n; i++ {
		start := rand.Int63()
		res[i] = [][]byte{
			[]byte(fmt.Sprintf("%016x", start)),
			[]byte(fmt.Sprintf("%016x", start+rand.Int63n(math.MaxInt64-start))),
		}
	}
	return res
}

func BenchmarkFind(b *testing.B) {
	for _, n := range []int{1 << 8, 1 << 12, 1 << 16} {
		tree := New[*int]()
		keys := testGenKeys(1 << 16)
		for j, k := range keys {
			tree.Insert(k[0], k[1], &j)
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for range b.N {
				k := testGenKey()
				tree.Find(k)
			}
		})
	}
}

func testGenKey() []byte {
	return []byte(fmt.Sprintf("%016x", rand.Int63()))
}

func BenchmarkFindAny(b *testing.B) {
	for _, n := range []int{1 << 8, 1 << 12, 1 << 16} {
		tree := New[*int]()
		keys := testGenKeys(n)
		for j, k := range keys {
			tree.Insert(k[0], k[1], &j)
		}
		b.ResetTimer()
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			for range b.N {
				keys := testGenKeyAny(16)
				tree.FindAny(keys...)
			}
		})
	}
}

func testGenKeyAny(n int) (keys [][]byte) {
	for range n {
		keys = append(keys, []byte(fmt.Sprintf("%016x", rand.Int63())))
	}
	return
}
