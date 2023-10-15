package html

func NewScript(text string) *Tag {
	script := NewTag("script")
	script.AddAttribute(NewAttribute("type", "text/javascript"))
	script.AddElements(NewText(text))
	return script
}

func NewDiv() *Tag {
	return NewTag("div")
}

func NewSpan(text string) *Tag {
	span := NewTag("span")
	span.AddElements(NewText(text))
	return span
}
