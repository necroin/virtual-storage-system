package runner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

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

	openParh := openRequest.Path
	selfUrl := fmt.Sprintf("%s/%s", runner.config.Url, runner.config.User.Token)
	if selfUrl != openRequest.SrcUrl {
		openParh = fmt.Sprintf("./tmp/%s", filepath.Base(openRequest.Path))

		copyRequest := &connector.CopyRequest{
			OldPath: openRequest.Path,
			NewPath: openParh,
			SrcUrl:  openRequest.SrcUrl,
		}
		connector.SendPostRequest(fmt.Sprintf("%s/storage/copy/%s", selfUrl, openRequest.Type), copyRequest)
	}

	logger.Info("[Runner] [OpenFileHandler] open %s", openParh)

	execTool, execArgs := runner.GetRunCommand(openParh)

	go func() {
		cmd := exec.Command(execTool, execArgs...)
		if err := cmd.Run(); err != nil {
			openResponse.Error = err
			logger.Error("[Runner] [OpenFileHandler] failed start process: %s", err)
			return
		}
	}()
}
