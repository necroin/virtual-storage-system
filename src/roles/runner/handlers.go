package runner

import (
	"encoding/json"
	"net/http"
	"os/exec"

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
	openResponse := &connector.OpenResponse{}
	defer json.NewEncoder(responseWriter).Encode(openResponse)

	openRequest := &connector.OpenRequest{}
	if err := json.NewDecoder(request.Body).Decode(openRequest); err != nil {
		openResponse.Error = err
		logger.Error("[Runner] [OpenFileHandler] failed decode message: %s", err)
		return
	}

	// openRequest.Path = filepath.Join(strings.Split(openRequest.Path, "/")...)
	logger.Info("[Runner] [OpenFileHandler] open %s", openRequest.Path)

	execTool, execArgs := runner.GetRunCommand(openRequest.Path)

	go func() {
		cmd := exec.Command(execTool, execArgs...)
		if err := cmd.Run(); err != nil {
			openResponse.Error = err
			logger.Error("[Runner] [OpenFileHandler] failed start process: %s", err)
			return
		}
	}()
}
