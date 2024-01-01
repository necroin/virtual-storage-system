package router

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	_ "embed"

	"vss/src/connector"
	"vss/src/roles"
	"vss/src/settings"
)

var (
	//go:embed assets/topology.html
	topologyHandlerResponseTemplate string
)

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	message := &connector.NotifyMessage{}
	if err := json.NewDecoder(request.Body).Decode(message); err != nil {
		// TODO: handle error
	}

	message.Url = strings.Split(request.RemoteAddr, ":")[0] + settings.DefaultPort
	if message.Url == ("127.0.0.1" + settings.DefaultPort) {
		message.Url = "localhost" + settings.DefaultPort
	}

	if message.Type == connector.NotifyMessageStorageType {
		router.storages[message.Hostname] = *message
		router.hostnames[message.Hostname] = path.Join(message.Url, message.Token)
		router.NotifyRunners()
	}

	if message.Type == connector.NotifyMessageRunnerType {
		router.runners[message.Hostname] = *message
		router.NotifyRunner(*message)
	}
}

func (router *Router) ExplorerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.ExplorerHandler(router, responseWriter, request)
}

func (router *Router) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.FilesystemHandler(router, responseWriter, request)
}

func (router *Router) DevicesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.DevicesHandler(router, responseWriter, request)
}
