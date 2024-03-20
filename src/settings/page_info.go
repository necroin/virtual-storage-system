package settings

import (
	_ "embed"
	"fmt"
	"os"
)

var (
	//go:embed assets/icons/home.svg
	BarHomeIcon string
	//go:embed assets/icons/filesystem.svg
	BarFilesystemIcon string
	//go:embed assets/icons/settings.svg
	BarSettingsIcon string
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
	//go:embed assets/explorer/icons/options.svg
	ExplorerIconOptions string
	//go:embed assets/explorer/icons/arrow_left.svg
	ExplorerIconArrowLeft string
	//go:embed assets/explorer/icons/status_bar_success.svg
	ExplorerStatusBarSuccess string
	//go:embed assets/explorer/icons/status_bar_fail.svg
	ExplorerStatusBarFail string
)

func GetExplorerPage() string {
	data, _ := os.ReadFile("src/settings/assets/explorer/explorer.html")
	return string(data)
}

func GetExplorerPageStyle() string {
	data, _ := os.ReadFile("src/settings/assets/explorer/explorer.css")
	style := string(data)
	return fmt.Sprintf(`<style type="text/css">%s</style>`, style)
}

func GetExplorerPageScript() string {
	data, _ := os.ReadFile("src/settings/assets/explorer/script.js")
	script := string(data)
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, script)
}

func GetHomePage() string {
	data, _ := os.ReadFile("src/settings/assets/home/page.html")
	return string(data)
}

func GetHomePageStyle() string {
	data, _ := os.ReadFile("src/settings/assets/home/style.css")
	style := string(data)
	return fmt.Sprintf(`<style type="text/css">%s</style>`, style)
}

func GetHomePageScript() string {
	data, _ := os.ReadFile("src/settings/assets/home/script.js")
	script := string(data)
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, script)
}

func GetSettingsPage() string {
	data, _ := os.ReadFile("src/settings/assets/settings/page.html")
	return string(data)
}

func GetSettingsPageStyle() string {
	data, _ := os.ReadFile("src/settings/assets/settings/style.css")
	style := string(data)
	return fmt.Sprintf(`<style type="text/css">%s</style>`, style)
}

func GetSettingsPageScript() string {
	data, _ := os.ReadFile("src/settings/assets/settings/script.js")
	script := string(data)
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, script)
}
