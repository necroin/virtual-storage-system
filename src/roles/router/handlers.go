package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"slices"
	"strings"
	"text/template"

	_ "embed"

	"vss/src/logger"
	"vss/src/message"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"
)

func (router *Router) GetTopologyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	response := &message.Topology{}
	for _, storage := range router.storages {
		response.Storages = append(response.Storages, storage)
	}
	for _, runner := range router.runners {
		response.Runners = append(response.Runners, runner)
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func (router *Router) NotifyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	notifyMessage := &message.Notify{}
	if err := json.NewDecoder(request.Body).Decode(notifyMessage); err != nil {
		logger.Error("[Router] [NotifyHandler] failed decode message: %s", err)
		return
	}

	notifyMessage.Url = strings.Split(request.RemoteAddr, ":")[0] + settings.DefaultPort
	if notifyMessage.Url == ("127.0.0.1" + settings.DefaultPort) {
		notifyMessage.Url = "localhost" + settings.DefaultPort
	}

	if notifyMessage.Type == message.NotifyMessageStorageType {
		router.storages[notifyMessage.Hostname] = *notifyMessage
		router.hostnames[notifyMessage.Hostname] = path.Join(notifyMessage.Url, notifyMessage.Token)
	}

	if notifyMessage.Type == message.NotifyMessageRunnerType {
		router.runners[notifyMessage.Hostname] = *notifyMessage
	}
}

func (router *Router) ExplorerHandler(responseWriter http.ResponseWriter, request *http.Request) {
	pageInfo := message.PageInfo{
		Url:               router.GetUrl(),
		Token:             router.config.User.Token,
		Style:             settings.GetExplorerPageStyle(),
		Script:            settings.GetExplorerPageScript(),
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
	openResponse := &message.OpenResponse{}
	defer json.NewEncoder(responseWriter).Encode(openResponse)

	openRequest := &message.OpenRequest{}
	if err := json.NewDecoder(request.Body).Decode(openRequest); err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[Router] [OpenFileHandler] failed decode message: %s", err)
		return
	}

	for _, runner := range router.runners {
		if runner.Platform == openRequest.Platform {
			logger.Info("[Router] [OpenFileHandler] send open request to %s runner on %s platform", runner.Hostname, runner.Platform)
			response, err := router.connector.SendPostRequest(runner.Url+utils.FormatTokemizedEndpoint(settings.RunnerOpenEndpoint, runner.Token), openRequest)
			if err != nil {
				logger.Error("[Router] [OpenFileHandler] selected runner failed execute: %s", err)
				continue
			}

			runnerOpenResponse := &message.OpenResponse{}
			if err := json.NewDecoder(response.Body).Decode(runnerOpenResponse); err != nil {
				logger.Error("[Router] [OpenFileHandler] failed decode open response: %s", err)
				continue
			}

			if runnerOpenResponse.Error != nil {
				logger.Error("[Router] [OpenFileHandler] %s runner failed execute: %s", runner.Hostname, err)
				continue
			}

			openResponse.Pid = runnerOpenResponse.Pid
			openResponse.RunnerUrl = path.Join(runner.Url, runner.Token)
			openResponse.Error = runnerOpenResponse.Error
			openResponse.StatusBar = message.StatusBarResponse{
				Status: settings.ExplorerStatusBarSuccess,
				Text:   fmt.Sprintf("File opened on runner %s", runner.Hostname),
			}
			return
		}
	}

	openResponse.StatusBar = message.StatusBarResponse{
		Status: settings.ExplorerStatusBarSuccess,
		Text:   "There is no way to run/open a file",
	}
}

func (router *Router) FiltersGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	json.NewEncoder(responseWriter).Encode(router.config.Settings.Filters)
}

func (router *Router) FiltersAddHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data, _ := io.ReadAll(request.Body)
	filter := string(data)
	if router.config.Settings.Filters.CurrentList == "Black list" {
		if !slices.Contains(router.config.Settings.Filters.BlackList, filter) {
			router.config.Settings.Filters.BlackList = append(router.config.Settings.Filters.BlackList, filter)
		}
	} else {
		if !slices.Contains(router.config.Settings.Filters.WhiteList, filter) {
			router.config.Settings.Filters.WhiteList = append(router.config.Settings.Filters.WhiteList, filter)
		}
	}
	router.config.Settings.Dump()
}

func (router *Router) FiltersRemoveHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data, _ := io.ReadAll(request.Body)
	filter := string(data)
	if router.config.Settings.Filters.CurrentList == "Black list" {
		router.config.Settings.Filters.BlackList = slices.DeleteFunc(router.config.Settings.Filters.BlackList, func(value string) bool { return value == filter })
	} else {
		router.config.Settings.Filters.WhiteList = slices.DeleteFunc(router.config.Settings.Filters.WhiteList, func(value string) bool { return value == filter })
	}
	router.config.Settings.Dump()
}

func (router *Router) FiltersSwapHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if router.config.Settings.Filters.CurrentList == "Black list" {
		router.config.Settings.Filters.CurrentList = "White list"
	} else {
		router.config.Settings.Filters.CurrentList = "Black list"
	}
	router.config.Settings.Dump()
}
