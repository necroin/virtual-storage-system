package runner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "embed"

	"vss/src/connector"
)

var (
	//go:embed assets/topology.html
	topologyHandlerResponseTemplate string
)

func (runner *Runner) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(fmt.Sprintf(topologyHandlerResponseTemplate, strings.Join(runner.storages, "</li>\n<li>"))))
}

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
