package html

import (
	"fmt"
	"strings"
)

type Head struct {
	elements []Element
}

func NewHead() *Head {
	return &Head{
		elements: []Element{},
	}
}

func (head *Head) ToHTML() string {
	items := []string{}
	for _, element := range head.elements {
		items = append(items, element.ToHTML())
	}
	return fmt.Sprintf("<head>%s</head>", strings.Join(items, ""))
}

func (head *Head) Add(elements ...Element) *Head {
	head.elements = append(head.elements, elements...)
	return head
}
