package html

import (
	"fmt"
	"strings"
)

type Form struct {
	elements []Element
}

func NewForm() *Form {
	return &Form{
		elements: []Element{},
	}
}

func (form *Form) ToHTML() string {
	items := []string{}
	for _, element := range form.elements {
		items = append(items, element.ToHTML())
	}
	return fmt.Sprintf("<form>%s</form>", strings.Join(items, ""))
}

func (form *Form) Add(elements ...Element) *Form {
	form.elements = append(form.elements, elements...)
	return form
}
