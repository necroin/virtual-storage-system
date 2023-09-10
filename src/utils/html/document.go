package html

type Document struct {
	*Tag
}

func NewDocument() *Document {
	return &Document{NewTag("html")}
}

func (document *Document) ToHTML() string {
	return "<!DOCTYPE html>" + document.Tag.ToHTML()
}
