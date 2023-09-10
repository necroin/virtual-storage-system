package html

func NewScript(text string) *Tag {
	script := NewTag("script")
	script.AddModifiers(NewModifier("type", "text/javascript"))
	script.AddElements(NewText(text))
	return script
}

func NewForm() *Tag {
	return NewTag("form")
}
