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
	r := New(nodes(), xxhash.Sum64String)
	t.Log(r.Lookup("hello"))
	t.Log(r.LookupN("hello", 3))
}

func TestLookupN_Lookup_Equivalence(t *testing.T) {
	r := New(nodes(), hashString)

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
	r := New(nodes(), xxhash.Sum64String)
	for i := 0; i < b.N; i++ {
		r.Lookup("github.com/foo/bar")
	}
}

func BenchmarkLookupN(b *testing.B) {
	b.ReportAllocs()
	r := New(nodes(), xxhash.Sum64String)
	for i := 0; i < b.N; i++ {
		r.LookupN("github.com/foo/bar", 3)
	}
}

func nodes() []string {
	nodes := make([]string, 24)
	for i := 0; i < cap(nodes); i++ {
		nodes[i] = fmt.Sprintf("indexed-search-%d.indexed-search:6070", i)
	}
	return nodes
}
