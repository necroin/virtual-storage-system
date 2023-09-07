package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vss/src/connector"
)

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	response := connector.TopologyMessage{
		Storages: router.storages,
		Runners:  router.runners,
	}

	json.NewEncoder(responseWriter).Encode(response)
}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	message := &connector.NotifyMessage{}
	if err := json.NewDecoder(request.Body).Decode(message); err != nil {
		// TODO: handle error
	}

	fmt.Println(message)
	if message.Type == connector.NotifyMessageStorageType {
		router.storages = append(router.storages, message.Url)
		router.NotifyRunners()
	}

	if message.Type == connector.NotifyMessageRunnerType {
		router.runners = append(router.runners, message.Url)
		router.NotifyRunner(message.Url)
	}
}

func (router *Router) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// TODO: ViewHandler
}
