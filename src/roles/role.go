package roles

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
	"vss/src/connector"
	"vss/src/settings"
)

type Role interface {
	CollectFileSystem(walkPath string) connector.FilesystemDirectory
	GetUrl() string
	GetHostnames() map[string]string
}

func FilesystemHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	fileSystemMessage := role.CollectFileSystem(string(msgPath))
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func DevicesHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	hostnames := role.GetHostnames()
	json.NewEncoder(responseWriter).Encode(hostnames)
}

func ExplorerHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	pageInfo := connector.PageInfo{
		Url:           role.GetUrl(),
		Style:         settings.GetExplorerStyle(),
		Script:        settings.GetExplorerScript(),
		IconCreate:    settings.ExplorerIconCreate,
		IconCut:       settings.ExplorerIconCut,
		IconCopy:      settings.ExplorerIconCopy,
		IconPaste:     settings.ExplorerIconPaste,
		IconDelete:    settings.ExplorerIconDelete,
		IconOptions:   settings.ExplorerIconOptions,
		IconArrowLeft: settings.ExplorerIconArrowLeft,
		StatusBarIcon: settings.ExplorerStatusBarSuccess,
	}
	pageTemplate, _ := template.New("ExplorerPage").Parse(settings.GetExplorerPage())
	pageTemplate.Execute(responseWriter, pageInfo)
}
