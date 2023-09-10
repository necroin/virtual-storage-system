package html

import (
	"fmt"
	"strings"
)

type Button struct {
	text    string
	onClick string
}

func NewButton(text string) *Button {
	return &Button{
		text: text,
	}
}

func (button *Button) ToHTML() string {
	modifiers := []string{}
	if button.onClick != "" {
		modifiers = append(modifiers, fmt.Sprintf(`onclick="%s"`, button.onClick))
	}
	return fmt.Sprintf("<button %s>%s</button>", strings.Join(modifiers, " "), button.text)
}

func (button *Button) SetOnClick(script string) *Button {
	button.onClick = script
	return button
}
