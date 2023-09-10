package html

import (
	"fmt"
	"strings"
)

type Button struct {
	*Tag
	icon    string
	text    string
	onClick string
}

func NewButton(text, icon string) *Button {
	return &Button{
		Tag:  NewTag("button"),
		icon: icon,
		text: text,
	}
}

func (button *Button) ToHTML() string {
	modifiers := []string{}
	if button.onClick != "" {
		modifiers = append(modifiers, fmt.Sprintf(`onclick="%s"`, button.onClick))
	}

	return fmt.Sprintf("<button %s>%s</button>", strings.Join(modifiers, " "), strings.Join([]string{button.icon, button.text}, " "))
}

func (button *Button) SetOnClick(script string) *Button {
	button.onClick = script
	return button
}

func (button *Button) SetIcon(name string) *Button {
	button.icon = name
	return button
}
