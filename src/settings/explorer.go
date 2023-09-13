package settings

import (
	_ "embed"
)

var (
	//go:embed assets/explorer/explorer.html
	ExplorerTemlate string
	//go:embed assets/explorer/explorer.css
	ExplorerStyle string
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
	//go:embed assets/explorer/script.js
	ExplorerScript string
)
