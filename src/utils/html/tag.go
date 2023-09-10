package html

import (
	"fmt"
	"strings"
)

type Element interface {
	ToHTML() string
}

type Modifier struct {
	name  string
	value string
}

func NewModifier(name string, value string) *Modifier {
	return &Modifier{
		name:  name,
		value: value,
	}
}

func (mod *Modifier) ToHTML() string {
	return fmt.Sprintf(`%s="%s"`, mod.name, mod.value)
}

type Tag struct {
	name      string
	elements  []Element
	modifiers []*Modifier
}

func NewTag(name string) *Tag {
	return &Tag{
		name:      name,
		elements:  []Element{},
		modifiers: []*Modifier{},
	}
}

func (tag *Tag) ToHTML() string {
	elements := []string{}
	for _, element := range tag.elements {
		elements = append(elements, element.ToHTML())
	}

	modifiers := []string{}
	for _, modifier := range tag.modifiers {
		modifiers = append(modifiers, modifier.ToHTML())
	}

	return fmt.Sprintf("<%s %s>%s</%s>", tag.name, strings.Join(modifiers, " "), strings.Join(elements, ""), tag.name)
}

func (tag *Tag) AddElements(elements ...Element) {
	tag.elements = append(tag.elements, elements...)
}

func (tag *Tag) AddModifiers(modifiers ...*Modifier) {
	tag.modifiers = append(tag.modifiers, modifiers...)
}
