package html

import (
	"strings"
)

type Button struct {
	*Tag
	icon    string
	text    string
	onClick string
}

func NewButton(text, icon string) *Button {
	button := &Button{
		Tag:  NewTag("button"),
		icon: icon,
		text: text,
	}
	button.AddElements(NewText(strings.Join([]string{button.icon, button.text}, " ")))
	return button
}

func (button *Button) SetOnClick(script string) *Button {
	button.Tag.AddAttribute(NewAttribute("onclick", script))
	return button
}

func (button *Button) SetIcon(name string) *Button {
	button.icon = name
	return button
}
