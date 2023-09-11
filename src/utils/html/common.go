package html

func NewScript(text string) *Tag {
	script := NewTag("script")
	script.AddAttribute(NewAttribute("type", "text/javascript"))
	script.AddElements(NewText(text))
	return script
}
