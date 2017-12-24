package rendezvous

type Rendezvous struct {
	nodes []string
	nhash []uint64
	hash  Hasher
}

type Hasher func(s string) uint64

func New(nodes []string, hash Hasher) *Rendezvous {
	r := &Rendezvous{
		nodes: make([]string, len(nodes)),
		nhash: make([]uint64, len(nodes)),
		hash:  hash,
	}

	for i, n := range nodes {
		r.nodes[i] = n
		r.nhash[i] = hash(n)
	}

	return r
}

func (r *Rendezvous) Lookup(k string) string {
	khash := r.hash(k)

	var midx int
	var mhash = khash ^ r.nhash[0]

	for i, nhash := range r.nhash[1:] {
		if h := khash ^ nhash; h < mhash {
			midx = i + 1
			mhash = h
		}
	}

	return r.nodes[midx]
}

func (r *Rendezvous) Add(node string) {
	r.nodes = append(r.nodes, node)
	r.nhash = append(r.nhash, r.hash(node))
}

func (r *Rendezvous) Remove(n int) {
	r.nodes = append(r.nodes[:n], r.nodes[n+1:]...)
	r.nhash = append(r.nhash[:n], r.nhash[n+1:]...)
}
