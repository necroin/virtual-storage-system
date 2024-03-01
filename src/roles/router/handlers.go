package router

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
	"text/template"

	_ "embed"

	"vss/src/connector"
	"vss/src/logger"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"
)

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	response := &connector.TopologyMessage{}
	for _, storage := range router.storages {
		response.Storages = append(response.Storages, storage)
	}
	for _, runner := range router.runners {
		response.Runners = append(response.Runners, runner)
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	message := &connector.NotifyMessage{}
	if err := json.NewDecoder(request.Body).Decode(message); err != nil {
		logger.Error("[Router] [NotifyHandler] failed decode message: %s", err)
		return
	}

	message.Url = strings.Split(request.RemoteAddr, ":")[0] + settings.DefaultPort
	if message.Url == ("127.0.0.1" + settings.DefaultPort) {
		message.Url = "localhost" + settings.DefaultPort
	}

	if message.Type == connector.NotifyMessageStorageType {
		router.storages[message.Hostname] = *message
		router.hostnames[message.Hostname] = path.Join(message.Url, message.Token)
		router.NotifyRunners()
	}

	if message.Type == connector.NotifyMessageRunnerType {
		router.runners[message.Hostname] = *message
		router.NotifyRunner(*message)
	}
}

func (router *Router) ExplorerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	pageInfo := connector.PageInfo{
		Url:               router.GetUrl(),
		Token:             router.config.User.Token,
		Style:             settings.GetExplorerStyle(),
		Script:            settings.GetExplorerScript(),
		IconCreate:        settings.ExplorerIconCreate,
		IconCut:           settings.ExplorerIconCut,
		IconCopy:          settings.ExplorerIconCopy,
		IconPaste:         settings.ExplorerIconPaste,
		IconDelete:        settings.ExplorerIconDelete,
		IconOptions:       settings.ExplorerIconOptions,
		IconArrowLeft:     settings.ExplorerIconArrowLeft,
		StatusBarIcon:     settings.ExplorerStatusBarSuccess,
		BarHomeIcon:       settings.BarHomeIcon,
		BarFilesystemIcon: settings.BarFilesystemIcon,
		BarSettingsIcon:   settings.BarSettingsIcon,
	}
	pageTemplate, _ := template.New("ExplorerPage").Parse(settings.GetExplorerPage())
	pageTemplate.Execute(responseWriter, pageInfo)
}

func (router *Router) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.FilesystemHandler(router, responseWriter, request)
}

func (router *Router) DevicesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.DevicesHandler(router, responseWriter, request)
}

func (router *Router) OpenFileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	openRequest := &connector.OpenRequest{}
	if err := json.NewDecoder(request.Body).Decode(openRequest); err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[Router] [OpenFileHandler] failed decode message: %s", err)
		return
	}

	for _, runner := range router.runners {
		if runner.Platform == openRequest.Platform {
			logger.Info("[Router] [OpenFileHandler] send open request to %s runner on %s platform", runner.Hostname, runner.Platform)
			response, err := connector.SendPostRequest(runner.Url+utils.FormatTokemizedEndpoint(settings.RunnerOpenEndpoint, runner.Token), openRequest)
			if err != nil {
				logger.Error("[Router] [OpenFileHandler] selected runner failed execute: %s", err)
				continue
			}

			openResponse := &connector.OpenResponse{}
			if err := json.NewDecoder(response.Body).Decode(openResponse); err != nil {
				logger.Error("[Router] [OpenFileHandler] failed decode open response: %s", err)
				continue
			}

			if openResponse.Error != nil {
				logger.Error("[Router] [OpenFileHandler] failed decode open response: %s", err)
				continue
			}

			roles.HandlerSuccess(responseWriter, openResponse.Message)
			break
		}
	}
	roles.HandlerSuccess(responseWriter, "Нет возможности запустить/открыть файл")
}
