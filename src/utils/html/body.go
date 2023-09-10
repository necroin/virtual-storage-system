package html

type Head struct {
	*Tag
}

func NewHead() *Head {
	return &Head{NewTag("head")}
}

type Body struct {
	*Tag
	head *Head
}

func NewBody(head *Head) *Body {
	return &Body{
		Tag:  NewTag("body"),
		head: head,
	}
}

func (body *Body) ToHTML() string {
	result := body.Tag.ToHTML()
	if body.head != nil {
		result = body.head.ToHTML() + result
	}
	return result
}
