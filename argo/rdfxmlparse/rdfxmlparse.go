package main

import (
	"fmt"
	"github.com/kierdavis/go/argo"
	"os"
	"strings"
)

var TestString = `
<rdf:RDF
    xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
    xmlns:foaf="http://xmlns.com/foaf/0.1/"
>
    <foaf:Person rdf:about="http://example.com/me">
        <foaf:name>Kier</foaf:name>
        <foaf:homepage rdf:resource="http://kierdavis.com/"/>
    </foaf:Person>
</rdf:RDF>
`

func main() {
	graph := argo.NewGraph()

	errChan := graph.Parse(argo.NewRdfXmlIO(), strings.NewReader(TestString))

	for err := range errChan {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	for triple := range graph.IterTriples() {
		fmt.Println(triple)
	}
}
