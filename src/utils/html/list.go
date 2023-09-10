package html

type UnorderedList struct {
	*Tag
}

func NewUnorderedList() *UnorderedList {
	return &UnorderedList{NewTag("ul")}
}

func (list *UnorderedList) AddElements(elements ...Element) *UnorderedList {
	for _, element := range elements {
		item := NewTag("li")
		item.AddElements(element)
		list.Tag.AddElements(item)
	}
	return list
}
