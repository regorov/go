package argo

import (
	"fmt"
)

// Generate thread-safe unique IDs

var newBlankIDChan = make(chan int)

func init() {
	go func() {
		i := 1

		for {
			newBlankIDChan <- i
			i++
		}
	}()
}

type Term interface {
	String() string
}

type Resource struct {
	URI string
}

func NewResource(uri string) (term Term) {
	return Term(&Resource{URI: uri})
}

func (term *Resource) String() (str string) {
	return fmt.Sprintf("<%s>", term.URI)
}

type Literal struct {
	Value		string
	Language	string
	Datatype	Term
}

func NewLiteral(value string) (term Term) {
	return Term(&Literal{Value: value})
}

func NewLiteralWithLanguage(value string, language string) (term Term) {
	return Term(&Literal{Value: value, Language: language})
}

func NewLiteralWithDatatype(value string, datatype Term) (term Term) {
	return Term(&Literal{Value: value, Datatype: datatype})
}

func (term *Literal) String() (str string) {
	str = fmt.Sprintf("\"%s\"", term.Value)

	if term.Language != "" {
		str += "@" + term.Language
	} else if term.Datatype != nil {
		str += "^^" + term.Datatype.String()
	}

	return str
}

type Blank struct {
	ID string
}

func NewBlank(id string) (term Term) {
	return Term(&Blank{ID: id})
}

func NewEmptyBlank() (term Term) {
	id := fmt.Sprintf("b%d", <-newBlankIDChan)
	return NewBlank(id)
}

func (term *Blank) String() (str string) {
	return "_:" + term.ID
}
