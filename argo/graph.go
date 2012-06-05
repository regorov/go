package argo

import (
	"strings"
)

type Graph struct {
	store		Store
	prefixes	map[string]string
}

func NewGraph(store Store) (graph *Graph) {
	graph = &Graph{store: store}
	graph.init()

	return graph
}

func NewGraphWithMemoryStore() (graph *Graph) {
	graph = &Graph{store: Store(NewMemoryStore())}
	graph.init()

	return graph
}

func (graph *Graph) init() {
	graph.prefixes = map[string]string{"http://www.w3.org/1999/02/22-rdf-syntax-ns#": "rdf"}
}

func (graph *Graph) Bind(uri string, prefix string) (ns Namespace) {
	graph.prefixes[uri] = prefix
	return NewNamespace(uri)
}

func (graph *Graph) LookupAndBind(prefix string) (ns Namespace, err error) {
	uri, err := Lookup(prefix)
	if err != nil {
		return ns, err
	}

	return graph.Bind(uri, prefix), nil
}

func (graph *Graph) Store() (store Store) {
	return graph.store
}

func (graph *Graph) SetStore(store Store) {
	graph.store = store
}

func (graph *Graph) Add(triple *Triple) {
	graph.store.Add(triple)
}

func (graph *Graph) AddTriple(subject Term, predicate Term, object Term) {
	graph.store.Add(NewTriple(subject, predicate, object))
}

func (graph *Graph) AddQuad(subject Term, predicate Term, object Term, context Term) {
	graph.store.Add(NewQuad(subject, predicate, object, context))
}

func (graph *Graph) Remove(triple *Triple) {
	graph.store.Remove(triple)
}

func (graph *Graph) RemoveTriple(subject Term, predicate Term, object Term) {
	graph.store.Remove(NewTriple(subject, predicate, object))
}

func (graph *Graph) RemoveQuad(subject Term, predicate Term, object Term, context Term) {
	graph.store.Remove(NewQuad(subject, predicate, object, context))
}

func (graph *Graph) Clear() {
	graph.store.Clear()
}

func (graph *Graph) NumTriples() (n int) {
	return graph.store.NumTriples()
}

func (graph *Graph) PumpTriples(ch chan *Triple) {
	graph.store.PumpTriples(ch)
}

func (graph *Graph) Triples() (triples []*Triple) {
	ch := make(chan *Triple)
	triples = make([]*Triple, graph.store.NumTriples())
	graph.store.PumpTriples(ch)

	i := 0
	for triple := range ch {
		triples[i] = triple
		i++
	}

	return triples
}

func (graph *Graph) TriplesBySubject() (subjects map[Term][]*Triple) {
	ch := make(chan *Triple)
	subjects = make(map[Term][]*Triple)
	graph.store.PumpTriples(ch)

	for triple := range ch {
		subjects[triple.Subject] = append(subjects[triple.Subject], triple)
	}

	return subjects
}

func splitPrefix(uri string) (base string, name string) {
	index := strings.LastIndex(uri, "#") + 1

	if index > 0 {
		return uri[:index], uri[index:]
	}

	index = strings.LastIndex(uri, "/") + 1

	if index > 0 {
		return uri[:index], uri[index:]
	}

	return "", uri
}
