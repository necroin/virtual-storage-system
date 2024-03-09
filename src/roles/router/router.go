package router

import (
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/logger"
	"vss/src/server"
	"vss/src/settings"
	"vss/src/utils"
)

type Router struct {
	config    *config.Config
	server    *server.Server
	storages  map[string]connector.NotifyMessage
	runners   map[string]connector.NotifyMessage
	hostnames map[string]string
}

func New(config *config.Config, server *server.Server) (*Router, error) {
	router := &Router{
		config:    config,
		server:    server,
		storages:  map[string]connector.NotifyMessage{},
		runners:   map[string]connector.NotifyMessage{},
		hostnames: map[string]string{},
	}

	go func() {
		time.Sleep(time.Second * 30)
		currentTime := time.Now().UnixNano()
		currentTimeSeconds := time.Duration(currentTime).Seconds()

		deleteStorages := []string{}
		deleteRunners := []string{}

		for hostname, storage := range router.storages {
			timestampSeconds := time.Duration(storage.Timestamp).Seconds()
			if currentTimeSeconds-timestampSeconds > settings.DefaultInstanceRemoveSeconds {
				deleteStorages = append(deleteStorages, hostname)
			}
		}
		for hostname, runner := range router.runners {
			timestampSeconds := time.Duration(runner.Timestamp).Seconds()
			if currentTimeSeconds-timestampSeconds > settings.DefaultInstanceRemoveSeconds {
				deleteRunners = append(deleteRunners, hostname)
			}
		}

		for _, hostname := range deleteStorages {
			delete(router.storages, hostname)
		}
		for _, hostname := range deleteRunners {
			delete(router.runners, hostname)
		}
	}()
	return router, nil
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
	logger.Debug("[Router] [CollectStorageFileSystem] collect on %s", url+utils.FormatTokemizedEndpoint(settings.StorageFilesystemEndpoint, token))

	result, err := connector.SendGetRequest[connector.FilesystemDirectory](
		url+utils.FormatTokemizedEndpoint(settings.StorageFilesystemEndpoint, token),
		[]byte(walkPath),
	)
	if err != nil {
		return connector.FilesystemDirectory{}
	}
	if result == nil {
		return connector.FilesystemDirectory{}
	}
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
	return router.config.Url + "/" + router.config.User.Token
}

func (router *Router) GetHostnames() map[string]string {
	return router.hostnames
}
