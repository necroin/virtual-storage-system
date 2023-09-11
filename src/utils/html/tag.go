package html

import (
	"fmt"
	"strings"
)

type Element interface {
	ToHTML() string
}

type Attribute struct {
	name  string
	value string
}

func NewAttribute(name string, value string) *Attribute {
	return &Attribute{
		name:  name,
		value: value,
	}
}

func (attribute *Attribute) ToHTML() string {
	return fmt.Sprintf(`%s="%s"`, attribute.name, attribute.value)
}

type Tag struct {
	name       string
	elements   []Element
	attributes map[string]*Attribute
}

func NewTag(name string) *Tag {
	return &Tag{
		name:       name,
		elements:   []Element{},
		attributes: map[string]*Attribute{},
	}
}

func (tag *Tag) ToHTML() string {
	elements := []string{}
	for _, element := range tag.elements {
		elements = append(elements, element.ToHTML())
	}

	attributes := []string{}
	for _, attribute := range tag.attributes {
		attributes = append(attributes, attribute.ToHTML())
	}

	return fmt.Sprintf("<%s %s>%s</%s>", tag.name, strings.Join(attributes, " "), strings.Join(elements, ""), tag.name)
}

func (tag *Tag) AddElements(elements ...Element) *Tag {
	tag.elements = append(tag.elements, elements...)
	return tag
}

func (tag *Tag) AddAttribute(attributes ...*Attribute) *Tag {
	for _, attribute := range attributes {
		tag.attributes[attribute.name] = attribute
	}
	return tag
}
