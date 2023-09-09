package html

import (
	"fmt"
	"strings"
)

type Body struct {
	head     *Head
	elements []Element
}

func NewBody(head *Head) *Body {
	return &Body{
		head:     head,
		elements: []Element{},
	}
}

func (body *Body) ToHTML() string {
	items := []string{}
	for _, element := range body.elements {
		items = append(items, element.ToHTML())
	}

	result := fmt.Sprintf("<body>%s</body>", strings.Join(items, ""))
	if body.head != nil {
		result = body.head.ToHTML() + result
	}

	return result
}

func (body *Body) Add(elements ...Element) *Body {
	body.elements = append(body.elements, elements...)
	return body
}
