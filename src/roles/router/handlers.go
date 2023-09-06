package router

import (
	"encoding/json"
	"net/http"
	"vss/src/settings"
)

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	response := settings.TopologyMessage{
		Storages: router.storages,
		Runners:  router.runners,
	}

	json.NewEncoder(responseWriter).Encode(response)
}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	message := &settings.NotifyMessage{}
	if err := json.NewDecoder(request.Body).Decode(message); err != nil {
		// TODO: handle error
	}

	if message.Type == settings.NotifyMessageStorageType {
		router.storages = append(router.storages, message.Url)
	}

	if message.Type == settings.NotifyMessageRunnerType {
		router.runners = append(router.runners, message.Url)
	}

	router.NotifyRunners()
}

func (router *Router) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// TODO: ViewHandler
}
