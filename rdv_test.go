package rendezvous

import (
	"fmt"
	"hash/fnv"
	"testing"
	"testing/quick"

	"github.com/cespare/xxhash/v2"
)

func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func TestEmpty(t *testing.T) {
	r := New([]string{}, hashString)
	r.Lookup("hello")
}

func TestLookupN(t *testing.T) {
	r := New(nodes(defaultNodeCount), xxhash.Sum64String)
	t.Log(r.Lookup("hello"))
	t.Log(r.LookupN("hello", 3))
}

func TestLookupN_Lookup_Equivalence(t *testing.T) {
	r := New(nodes(defaultNodeCount), hashString)

	prop := func(k string) bool {
		a := r.LookupN(k, 1)[0]
		b := r.Lookup(k)
		return a == b
	}

	if err := quick.Check(prop, &quick.Config{MaxCount: 1e6}); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkLookup(b *testing.B) {
	b.ReportAllocs()
	for _, nodeCount := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("NodeCount-%d", nodeCount), func(b *testing.B) {
			r := New(nodes(nodeCount), xxhash.Sum64String)
			for i := 0; i < b.N; i++ {
				r.Lookup("github.com/foo/bar")
			}
		})
	}
}

func BenchmarkLookupN(b *testing.B) {
	b.ReportAllocs()
	b.ReportAllocs()
	for _, nodeCount := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("NodeCount-%d", nodeCount), func(b *testing.B) {
			r := New(nodes(nodeCount), xxhash.Sum64String)
			for i := 0; i < b.N; i++ {
				r.LookupN("github.com/foo/bar", 3)
			}

		})
	}
}

const defaultNodeCount = 24

func nodes(n int) []string {
	nodes := make([]string, n)
	for i := 0; i < cap(nodes); i++ {
		nodes[i] = fmt.Sprintf("indexed-search-%d.indexed-search:6070", i)
	}
	return nodes
}
