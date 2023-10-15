package router

import (
	"vss/src/config"
	"vss/src/connector"
	"vss/src/settings"
	"vss/src/utils"
)

type Router struct {
	config    *config.Config
	storages  []connector.NotifyMessage
	runners   []connector.NotifyMessage
	hostnames map[string]string
}

func New(config *config.Config) (*Router, error) {
	return &Router{
		config:    config,
		storages:  []connector.NotifyMessage{},
		runners:   []connector.NotifyMessage{},
		hostnames: map[string]string{},
	}, nil
}

func (router *Router) NotifyRunner(instance connector.NotifyMessage) {
	// topology := connector.TopologyMessage{
	// 	Storages: router.storages,
	// 	Runners:  router.runners,
	// }
	// connector.SendPostRequest(url+settings.RunnerNotifyEndpoint, topology)
}

func (router *Router) NotifyRunners() {
	for _, runner := range router.runners {
		router.NotifyRunner(runner)
	}
}

func (router *Router) CollectStorageFileSystem(url string, token string, walkPath string) connector.FilesystemDirectory {
	result, _ := connector.SendGetRequest[connector.FilesystemDirectory](
		url+utils.FormatTokemizedEndpoint(settings.StorageFilesystemEndpoint, token),
		[]byte(walkPath),
	)
	return *result
}

func (router *Router) CollectFileSystem(walkPath string) connector.FilesystemDirectory {
	fileSystemDirectory := connector.FilesystemDirectory{
		Directories: map[string]connector.FileInfo{},
		Files:       map[string]connector.FileInfo{},
	}

	for _, storage := range router.storages {
		storageFilesystem := router.CollectStorageFileSystem(
			storage.Url,
			storage.Token,
			walkPath,
		)
		for directory, info := range storageFilesystem.Directories {
			fileSystemDirectory.Directories[directory] = info
		}

		for file, info := range storageFilesystem.Files {
			fileSystemDirectory.Files[file] = info
		}
	}
	return fileSystemDirectory
}

func (router *Router) GetUrl() string {
	return router.config.RouterUrl
}

func (router *Router) GetMainEndpoint() string {
	return utils.FormatTokemizedEndpoint(settings.RouterMainEndpoint, router.config.User.Token)
}

func (router *Router) GetHostnames() map[string]string {
	return router.hostnames
}
