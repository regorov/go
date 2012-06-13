package argo

import (
	"encoding/xml"
	"io"
)

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
	NodeID   string `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# nodeID,nodeID"`
	Datatype string `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# datatype,nodeID"`
	Text     string `xml:",chardata"`
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
		var document xmlDocument

		decoder := xml.NewDecoder(r)
		err := decoder.Decode(&document)
		if err != nil {
			errChan <- err
			return
		}

		for _, description := range document.Descriptions {

		}

		close(tripleChan)
		close(errChan)
	}()

	return tripleChan, errChan
}
