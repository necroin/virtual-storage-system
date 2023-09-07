package runner

import (
	"encoding/json"
	"net/http"
	"vss/src/connector"
)

func (runner *Runner) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	topology := &connector.TopologyMessage{}
	if err := json.NewDecoder(request.Body).Decode(topology); err != nil {
		// TODO: handle error
	}
	runner.storages = append(runner.storages, topology.Storages...)
}

func (runner *Runner) OpenFileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// TODO: OpenFileHandler
}

func (runner *Runner) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// TODO: ViewHandler
}
