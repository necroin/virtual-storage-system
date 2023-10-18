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

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// responseWriter.Write([]byte(fmt.Sprintf(topologyHandlerResponseTemplate, strings.Join(router.storages, "</li>\n<li>"), strings.Join(router.runners, "</li>\n<li>"))))
}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	message := &connector.NotifyMessage{}
	if err := json.NewDecoder(request.Body).Decode(message); err != nil {
		// TODO: handle error
	}

	if message.Type == connector.NotifyMessageStorageType {
		message.Url = strings.Split(request.RemoteAddr, ":")[0] + settings.DefaultPort
		router.storages = append(router.storages, *message)
		router.hostnames[message.Hostname] = path.Join(message.Url, message.Token)
		router.NotifyRunners()
	}

	if message.Type == connector.NotifyMessageRunnerType {
		router.runners = append(router.runners, *message)
		router.NotifyRunner(*message)
	}
}

func (router *Router) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (router *Router) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.MainHandler(router, responseWriter, request)
}
