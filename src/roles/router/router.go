package router

import (
	"io/fs"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/settings"
)

type Router struct {
	url       string
	storages  []string
	runners   []string
	hostnames map[string]string
}

func New(config *config.Config) (*Router, error) {
	return &Router{
		url:       config.RouterUrl,
		hostnames: map[string]string{},
	}, nil
}

func (router *Router) NotifyRunner(url string) {
	topology := connector.TopologyMessage{
		Storages: router.storages,
		Runners:  router.runners,
	}
	connector.SendPostRequest(url+settings.RunnerNotifyEndpoint, topology)
}

func (router *Router) NotifyRunners() {
	for _, runner := range router.runners {
		router.NotifyRunner(runner)
	}
}

func (router *Router) CollectFileSystem(walkPath string) connector.FilesystemDirectory {
	fileSystemDirectory := connector.FilesystemDirectory{
		Directories: map[string]fs.FileInfo{},
		Files:       map[string]fs.FileInfo{},
	}
	return fileSystemDirectory
}

func (router *Router) GetUrl() string {
	return router.url
}

func (router *Router) GetMainEndpoint() string {
	return settings.RouterMainEndpoint
}

func (router *Router) GetHostnames() map[string]string {
	return router.hostnames
}
