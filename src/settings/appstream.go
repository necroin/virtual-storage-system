package settings

import (
	"fmt"
	"os"
)

func GetAppStreamPage() string {
	data, _ := os.ReadFile("src/settings/assets/appstream/appstream.html")
	return string(data)
}

func GetAppStreamPageStyle() string {
	data, _ := os.ReadFile("src/settings/assets/appstream/appstream.css")
	style := string(data)
	return fmt.Sprintf(`<style type="text/css">%s</style>`, style)
}

func GetAppStreamPageScript() string {
	data, _ := os.ReadFile("src/settings/assets/appstream/appstream.js")
	script := string(data)
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, script)
}
