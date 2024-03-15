package runner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"

	"vss/src/logger"
	"vss/src/message"
)

func (runner *Runner) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	topology := &message.TopologyMessage{}
	if err := json.NewDecoder(request.Body).Decode(topology); err != nil {
		logger.Error("[Runner] [NotifyHandler] failed decode message: %s", err)
		return
	}
	runner.storages = append(runner.storages, topology.Storages...)
}

func (runner *Runner) OpenFileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	openResponse := &message.OpenResponse{}
	defer json.NewEncoder(responseWriter).Encode(openResponse)

	openRequest := &message.OpenRequest{}
	if err := json.NewDecoder(request.Body).Decode(openRequest); err != nil {
		openResponse.Error = err
		logger.Error("[Runner] [OpenFileHandler] failed decode message: %s", err)
		return
	}

	openPath := openRequest.Path
	selfUrl := fmt.Sprintf("%s/%s", runner.config.Url, runner.config.User.Token)
	if selfUrl != openRequest.SrcUrl {
		openPath = fmt.Sprintf("./tmp/%s", filepath.Base(openRequest.Path))

		copyRequest := &message.CopyRequest{
			OldPath: openRequest.Path,
			NewPath: openPath,
			SrcUrl:  openRequest.SrcUrl,
		}
		runner.connector.SendPostRequest(fmt.Sprintf("%s/storage/copy/%s", selfUrl, openRequest.Type), copyRequest)
	}

	logger.Info("[Runner] [OpenFileHandler] open %s", openPath)

	execTool, execArgs := runner.GetRunCommand(openPath)

	go func() {
		cmd := exec.Command(execTool, execArgs...)
		if err := cmd.Run(); err != nil {
			openResponse.Error = err
			logger.Error("[Runner] [OpenFileHandler] failed start process: %s", err)
			return
		}
	}()
}
