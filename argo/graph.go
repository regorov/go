package argo

import (
	"io"
	"strings"
)

type Graph struct {
	triples  []*Triple
	prefixes map[string]string
}

func NewGraph() (graph *Graph) {
	return &Graph{
		triples:  make([]*Triple, 0),
		prefixes: map[string]string{"http://www.w3.org/1999/02/22-rdf-syntax-ns#": "rdf"},
	}
}

func (graph *Graph) Bind(uri string, prefix string) (ns Namespace) {
	graph.prefixes[uri] = prefix
	return NewNamespace(uri)
}

func (graph *Graph) LookupAndBind(prefix string) (ns Namespace, err error) {
	uri, err := LookupPrefix(prefix)
	if err != nil {
		return ns, err
	}

	return graph.Bind(uri, prefix), nil
}

func (graph *Graph) Add(triple *Triple) (index int) {
	index = len(graph.triples)
	graph.triples = append(graph.triples, triple)
	return index
}

func (graph *Graph) AddTriple(subject Term, predicate Term, object Term) {
	graph.Add(NewTriple(subject, predicate, object))
}

func (graph *Graph) AddQuad(subject Term, predicate Term, object Term, context Term) {
	graph.Add(NewQuad(subject, predicate, object, context))
}

func (graph *Graph) Remove(triple *Triple) {
	for i, t := range graph.triples {
		if t == triple {
			graph.RemoveIndex(i)
			return
		}
	}
}

func (graph *Graph) RemoveIndex(index int) {
	graph.triples = append(graph.triples[:index], graph.triples[index+1:]...)
}

func (graph *Graph) RemoveTriple(subject Term, predicate Term, object Term) {
	graph.Remove(NewTriple(subject, predicate, object))
}

func (graph *Graph) RemoveQuad(subject Term, predicate Term, object Term, context Term) {
	graph.Remove(NewQuad(subject, predicate, object, context))
}

func (graph *Graph) Clear() {
	graph.triples = graph.triples[:0]
}

func (graph *Graph) Num() (n int) {
	return len(graph.triples)
}

func (graph *Graph) IterTriples() (ch chan *Triple) {
	ch = make(chan *Triple)

	go func() {
		for _, triple := range graph.triples {
			ch <- triple
		}

		close(ch)
	}()

	return ch
}

func (graph *Graph) Triples() (triples []*Triple) {
	return graph.triples
}

func (graph *Graph) TriplesBySubject() (subjects map[Term][]*Triple) {
	subjects = make(map[Term][]*Triple)

	for triple := range graph.IterTriples() {
		subjects[triple.Subject] = append(subjects[triple.Subject], triple)
	}

	return subjects
}

func (graph *Graph) Parse(parser Parser, r io.Reader) (errChan chan error) {
	tripleChan, errChan := parser.Parse(r)

	go func() {
		for triple := range tripleChan {
			graph.Add(triple)
		}
	}()

	return errChan
}

func (graph *Graph) Serialize(serializer Serializer, w io.Writer) (err error) {
	return serializer.Serialize(w, graph.IterTriples())
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
