package settings

import (
	"fmt"
	"os"
)

func GetAuthenticationTemlate() string {
	data, _ := os.ReadFile("src/settings/assets/authentication/authentication.html")
	return string(data)
}

func GetAuthenticationStyle() string {
	data, _ := os.ReadFile("src/settings/assets/authentication/authentication.css")
	style := string(data)
	return fmt.Sprintf(`<style type="text/css">%s</style>`, style)
}

func GetAuthenticationScript() string {
	data, _ := os.ReadFile("src/settings/assets/authentication/authentication.js")
	script := string(data)
	return fmt.Sprintf(`<script type="text/javascript">%s</script>`, script)
}
