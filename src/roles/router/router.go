package router

import (
	"encoding/json"
	"fmt"
	"path"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/logger"
	"vss/src/message"
	"vss/src/settings"
	"vss/src/utils"

	"github.com/robfig/cron"
)

type Router struct {
	config           *config.Config
	connector        *connector.Connector
	storages         map[string]message.Notify
	runners          map[string]message.Notify
	hostnames        map[string]string
	replicationTasks map[string]*cron.Cron
}

func New(config *config.Config, connector *connector.Connector) (*Router, error) {
	router := &Router{
		config:           config,
		connector:        connector,
		storages:         map[string]message.Notify{},
		runners:          map[string]message.Notify{},
		hostnames:        map[string]string{},
		replicationTasks: map[string]*cron.Cron{},
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
			delete(router.hostnames, hostname)
		}
		for _, hostname := range deleteRunners {
			delete(router.runners, hostname)
		}
	}()

	return router, nil
}

func (router *Router) StartReplicationTasks() {
	for _, replicationCron := range router.replicationTasks {
		replicationCron.Start()
	}
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

func (router *Router) CollectFilesystem(walkPath string) message.FilesystemDirectory {
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

func (router *Router) SendOpenRequest(runner message.Notify, openRequest *message.OpenRequest) (*message.OpenResponse, error) {
	response, err := router.connector.SendPostRequest(runner.Url+utils.FormatTokemizedEndpoint(settings.RunnerOpenEndpoint, runner.Token), openRequest)
	if err != nil {
		return nil, fmt.Errorf("[Router] [SendOpenRequest] selected runner failed execute: %s", err)
	}

	openResponse := &message.OpenResponse{}
	if err := json.NewDecoder(response.Body).Decode(openResponse); err != nil {
		return nil, fmt.Errorf("[Router] [SendOpenRequest] failed decode open response: %s", err)
	}

	if openResponse.Error != nil {
		return nil, fmt.Errorf("[Router] [SendOpenRequest] %s runner failed execute: %s", runner.Hostname, err)
	}

	return openResponse, nil
}

func (router *Router) AddReplicationTask(replication config.ReplicationSettings) *cron.Cron {
	replicationCron := cron.New()

	replicationCron.AddFunc(replication.Cron, func() {
		srcStorage, ok := router.storages[replication.SrcHostname]
		if !ok {
			logger.Error("[Router] [Replcation] [Cron] failed source storage %s", replication.SrcHostname)
		}
		dstStorage, ok := router.storages[replication.DstHostname]
		if !ok {
			logger.Error("[Router] [Replcation] [Cron] failed destination storage %s", replication.DstHostname)
		}

		copyRequest := &message.CopyRequest{
			SrcPath: replication.SrcPath,
			DstPath: replication.DstPath,
			SrcUrl:  path.Join(srcStorage.Url, srcStorage.Token),
		}

		_, err := router.connector.SendPostRequest(dstStorage.Url+utils.FormatTokemizedEndpoint(settings.StorageCopyEndpoint, dstStorage.Token), copyRequest)
		if err != nil {
			logger.Error("[Router] [Replcation] [Cron] failed copy: %s", err)
		}
	})

	router.replicationTasks[replication.String()] = replicationCron
	return replicationCron
}

func (router *Router) AddReplicationTasks(replication ...config.ReplicationSettings) {
	for _, replication := range router.config.Settings.Replication {
		router.AddReplicationTask(replication)
	}
}
