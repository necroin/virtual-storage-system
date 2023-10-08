package settings

import (
	_ "embed"
	"os"
)

var (
	//go:embed assets/explorer/explorer.html
	ExplorerTemlate string
	//go:embed assets/explorer/explorer.css
	ExplorerStyle string
	//go:embed assets/explorer/script.js
	ExplorerScript string
	//go:embed assets/explorer/icons/create_plus.svg
	ExplorerIconCreate string
	//go:embed assets/explorer/icons/cut.svg
	ExplorerIconCut string
	//go:embed assets/explorer/icons/copy.svg
	ExplorerIconCopy string
	//go:embed assets/explorer/icons/paste.svg
	ExplorerIconPaste string
	//go:embed assets/explorer/icons/delete.svg
	ExplorerIconDelete string
	//go:embed assets/explorer/icons/arrow_left.svg
	ExplorerIconArrowLeft string
)

func GetExplorerTemlate() string {
	data, _ := os.ReadFile("assets/explorer/explorer.html")
	return string(data)
}

func GetExplorerStyle() string {
	data, _ := os.ReadFile("assets/explorer/explorer.css")
	return string(data)
}

func GetExplorerScript() string {
	data, _ := os.ReadFile("assets/explorer/script.js")
	return string(data)
}
