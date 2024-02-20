package runner

import (
	"encoding/json"
	"net/http"

	_ "embed"

	"vss/src/connector"
	"vss/src/logger"
)

func (runner *Runner) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	topology := &connector.TopologyMessage{}
	if err := json.NewDecoder(request.Body).Decode(topology); err != nil {
		logger.Error("[Runner] [NotifyHandler] failed decode message: %s", err)
		return
	}
	runner.storages = append(runner.storages, topology.Storages...)
}

func (runner *Runner) OpenFileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// TODO: OpenFileHandler
}
