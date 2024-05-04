package router

import (
	"fmt"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/logger"
	"vss/src/message"
	"vss/src/server"
	"vss/src/settings"
	"vss/src/utils"

	"gopkg.in/ini.v1"
)

type Router struct {
	config    *config.Config
	connector *connector.Connector
	server    *server.Server
	storages  map[string]message.Notify
	runners   map[string]message.Notify
	hostnames map[string]string
	filters   *ini.File
}

func New(config *config.Config, server *server.Server, connector *connector.Connector) (*Router, error) {
	filters, err := ini.LooseLoad("filters.ini")
	if err != nil {
		return nil, fmt.Errorf("[Router] failed create filters file")
	}

	router := &Router{
		config:    config,
		server:    server,
		connector: connector,
		storages:  map[string]message.Notify{},
		runners:   map[string]message.Notify{},
		hostnames: map[string]string{},
		filters:   filters,
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

func (router *Router) CollectStorageFileSystem(url string, token string, walkPath string) message.FilesystemDirectory {
	logger.Debug("[Router] [CollectStorageFileSystem] collect on %s", url+utils.FormatTokemizedEndpoint(settings.StorageFilesystemEndpoint, token))

	result := &message.FilesystemDirectory{}
	if err := router.connector.SendGetRequest(
		url+utils.FormatTokemizedEndpoint(settings.StorageFilesystemEndpoint, token),
		[]byte(walkPath),
		result,
	); err != nil {
		return message.FilesystemDirectory{}
	}

	return *result
}

func (router *Router) CollectFileSystem(walkPath string) message.FilesystemDirectory {
	fileSystemDirectory := message.FilesystemDirectory{
		Directories: map[string]message.FileInfo{},
		Files:       map[string][]message.FileInfo{},
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
			fileSystemDirectory.Files[file] = append(fileSystemDirectory.Files[file], info...)
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
