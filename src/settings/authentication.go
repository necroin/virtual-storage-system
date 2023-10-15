package settings

import "os"

var ()

func GetAuthenticationTemlate() string {
	data, _ := os.ReadFile("src/settings/assets/authentication/authentication.html")
	return string(data)
}

func GetAuthenticationStyle() string {
	data, _ := os.ReadFile("src/settings/assets/authentication/authentication.css")
	return string(data)
}

func GetAuthenticationScript() string {
	data, _ := os.ReadFile("src/settings/assets/authentication/authentication.js")
	return string(data)
}
