package roles

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"vss/src/connector"
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
