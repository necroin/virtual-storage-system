package buffer

var (
	text string
)

func SetText(value string) {
	text = value
}

func GetText() string {
	return text
}
