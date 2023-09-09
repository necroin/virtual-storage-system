package html

type Text struct {
	value string
}

func NewText(value string) *Text {
	return &Text{
		value: value,
	}
}

func (s *Text) ToHTML() string {
	return s.value
}
