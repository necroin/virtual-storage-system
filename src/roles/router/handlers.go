package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "embed"

	"vss/src/connector"
	"vss/src/roles"
)

var (
	//go:embed assets/topology.html
	topologyHandlerResponseTemplate string
)

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(fmt.Sprintf(topologyHandlerResponseTemplate, strings.Join(router.storages, "</li>\n<li>"), strings.Join(router.runners, "</li>\n<li>"))))
}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	message := &connector.NotifyMessage{}
	if err := json.NewDecoder(request.Body).Decode(message); err != nil {
		// TODO: handle error
	}

	if message.Type == connector.NotifyMessageStorageType {
		router.storages = append(router.storages, message.Url)
		router.hostnames[message.Hostname] = message.Url
		router.NotifyRunners()
	}

	if message.Type == connector.NotifyMessageRunnerType {
		router.runners = append(router.runners, message.Url)
		router.NotifyRunner(message.Url)
	}
}

func (router *Router) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.MainHandler(router, responseWriter, request)
}
