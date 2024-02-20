package settings

import (
	"fmt"
	"os"
)

func GetHomeTemplate() string {
	data, _ := os.ReadFile("src/settings/assets/home/page.html")
	return string(data)
}

func GetHomeStyle() string {
	data, _ := os.ReadFile("src/settings/assets/home/style.css")
	style := string(data)
	return fmt.Sprintf(`<style type="text/css">%s</style>`, style)
}

func GetHomeScript() string {
	data, _ := os.ReadFile("src/settings/assets/home/script.js")
	script := string(data)
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, script)
}
