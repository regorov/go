package argo

import (
	"encoding/xml"
	"fmt"
	"io"
)

/*
type xmlDocument struct {
	XMLName      xml.Name         `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# RDF"`
	Descriptions []xmlDescription `xml:",any"`
}

type xmlDescription struct {
	XMLName    xml.Name
	About      string        `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# about,attr"`
	NodeID     string        `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# nodeID,attr"`
	Properties []xmlProperty `xml:",any"`
}

type xmlProperty struct {
	XMLName  xml.Name
	Resource string `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# resource,attr"`
	NodeID   string `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# nodeID,attr"`
	Datatype string `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# datatype,attr"`
	Text     string `xml:",chardata"`
}
*/

const (
	stateTop = iota
	stateDescriptions
	stateProperties
	statePropertyValue
)

var (
	RdfNs = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"

	RdfRdf         = xml.Name{RdfNs, "RDF"}
	RdfDescription = xml.Name{RdfNs, "Description"}
	RdfAbout       = xml.Name{RdfNs, "about"}
	RdfNodeID      = xml.Name{RdfNs, "nodeID"}
	RdfResource    = xml.Name{RdfNs, "resource"}
	RdfDatatype    = xml.Name{RdfNs, "datatype"}

	XmlLang = xml.Name{"xml", "lang"}
)

func Name2Term(name xml.Name) (term Term) {
	return NewResource(name.Space + name.Local)
}

type RdfXmlIO struct {
}

func NewRdfXmlIO() (ps *RdfXmlIO) {
	return &RdfXmlIO{}
}

func (ps *RdfXmlIO) Parse(r io.Reader) (tripleChan chan *Triple, errChan chan error) {
	tripleChan = make(chan *Triple)
	errChan = make(chan error)

	go func() {
		decoder := xml.NewDecoder(r)
		state := stateTop

		var subject, predicate, datatype Term
		var language string

	loop:
		for {
			itok, err := decoder.Token()
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}

				break loop
			}

			switch state {
			case stateTop:
				switch tok := itok.(type) {
				case xml.StartElement:
					if tok.Name != RdfRdf {
						errChan <- fmt.Errorf("Syntax error: expected <rdf:RDF>")
						break loop
					}

					state = stateDescriptions
				}

			case stateDescriptions:
				switch tok := itok.(type) {
				case xml.StartElement:
					subject = nil
					extraAttrs := make([]xml.Attr, 0)

					for _, attr := range tok.Attr {
						if attr.Name == RdfAbout {
							subject = NewResource(attr.Value)
						} else if attr.Name == RdfNodeID {
							subject = NewNode(attr.Value)
						} else {
							extraAttrs = append(extraAttrs, attr)
						}
					}

					if subject == nil {
						subject = NewBlankNode()
					}

					if tok.Name != RdfDescription {
						tripleChan <- NewTriple(subject, A, Name2Term(tok.Name))
					}

					for _, attr := range extraAttrs {
						tripleChan <- NewTriple(subject, Name2Term(attr.Name), NewLiteral(attr.Value))
					}

					state = stateProperties

				case xml.EndElement: // Must be the toplevel tag (</rdf:RDF>)
					break loop
				}

			case stateProperties:
				switch tok := itok.(type) {
				case xml.StartElement:
					predicate = Name2Term(tok.Name)
					language = ""
					datatype = nil

					for _, attr := range tok.Attr {
						if attr.Name == RdfResource {
							tripleChan <- NewTriple(subject, predicate, NewResource(attr.Value))
							continue loop

						} else if attr.Name == RdfNodeID {
							tripleChan <- NewTriple(subject, predicate, NewNode(attr.Value))
							continue loop

						} else if attr.Name == RdfDatatype {
							datatype = NewResource(attr.Value)

						} else if attr.Name == XmlLang {
							language = attr.Value

						} else {
							errChan <- fmt.Errorf("Invalid attribute on property tag: %s:%s", attr.Name.Space, attr.Name.Local)
							break loop
						}
					}

					state = statePropertyValue

				case xml.EndElement: // Must be a description tag (</rdf:Description>)
					state = stateDescriptions
				}

			case statePropertyValue:
				switch tok := itok.(type) {
				case xml.CharData:
					tripleChan <- NewTriple(subject, predicate, NewLiteralWithLanguageAndDatatype(string(tok), language, datatype))

				case xml.EndElement: // Must be a property tag (</foaf:name>)
					state = stateProperties
				}
			}
		}

		close(tripleChan)
		close(errChan)
		return
	}()

	return tripleChan, errChan
}
