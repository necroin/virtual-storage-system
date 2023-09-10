package html

import (
	"fmt"
	"strings"
)

type UnorderedList struct {
	elements []Element
}

func NewUnorderedList() *UnorderedList {
	return &UnorderedList{
		elements: []Element{},
	}
}

func (list *UnorderedList) ToHTML() string {
	items := []string{}
	for _, element := range list.elements {
		items = append(items, fmt.Sprintf("<li>%v</li>", element.ToHTML()))
	}
	return fmt.Sprintf("<ul>%s</ul>", strings.Join(items, ""))
}

func (list *UnorderedList) Add(elements ...Element) *UnorderedList {
	list.elements = append(list.elements, elements...)
	return list
}
