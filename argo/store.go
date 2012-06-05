package argo

import ()

type Store interface {
	Add(*Triple)
	Remove(*Triple)
	Clear()
	NumTriples() int
	PumpTriples(chan *Triple)
}

type MemoryStore struct {
	triples []*Triple
}

func NewMemoryStore() (store *MemoryStore) {
	store = &MemoryStore{triples: make([]*Triple, 0)}
	return store
}

func (store *MemoryStore) Add(triple *Triple) {
	store.triples = append(store.triples, triple)
}

func (store *MemoryStore) Remove(triple *Triple) {
	for i, t := range store.triples {
		if t == triple {
			store.triples = append(store.triples[:i], store.triples[i+1:]...)
			return
		}
	}
}

func (store *MemoryStore) Clear() {
	store.triples = make([]*Triple, 0)
}

func (store *MemoryStore) NumTriples() (n int) {
	return len(store.triples)
}

func (store *MemoryStore) PumpTriples(ch chan *Triple) {
	go func() {
		for _, triple := range store.triples {
			ch <- triple
		}

		close(ch)
	}()
}
