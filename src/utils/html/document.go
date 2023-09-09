package html

import (
	"fmt"
	"strings"
)

type Element interface {
	ToHTML() string
}

type Document struct {
	elements []Element
}

func NewDocument() *Document {
	return &Document{
		elements: []Element{},
	}
}

func (document *Document) ToHTML() string {
	items := []string{}
	for _, element := range document.elements {
		items = append(items, element.ToHTML())
	}
	return fmt.Sprintf("<!DOCTYPE html><html>%s</html>", strings.Join(items, ""))
}

func (document *Document) Add(elements ...Element) *Document {
	document.elements = append(document.elements, elements...)
	return document
}
