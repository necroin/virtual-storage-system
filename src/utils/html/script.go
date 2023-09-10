package html

import "fmt"

type Script struct {
	text string
}

func NewScript(text string) *Script {
	return &Script{
		text: text,
	}
}

func (script *Script) ToHTML() string {
	return fmt.Sprintf("<script>%s</script>", script.text)
}
