package storage

import (
	"os"
	"path"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/message"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"
)

type Storage struct {
	config    *config.Config
	connector *connector.Connector
}

func New(config *config.Config, connector *connector.Connector) (*Storage, error) {

	return &Storage{
		config:    config,
		connector: connector,
	}, nil
}

func (storage *Storage) NotifyRouter(url string) error {
	token, err := roles.GetRouterToken(storage.connector, url, storage.config.User.Username, storage.config.User.Password)
	if err != nil {
		return err
	}

	message := message.Notify{
		Type:      message.NotifyMessageStorageType,
		Url:       storage.config.Url,
		Hostname:  storage.config.Hostname,
		Token:     storage.config.User.Token,
		Platform:  storage.config.Roles.Runner.Platform,
		Timestamp: time.Now().UnixNano(),
	}
	_, err = storage.connector.SendPostRequest(
		url+utils.FormatTokemizedEndpoint(settings.RouterNotifyEndpoint, token),
		message,
	)
	return err
}

func (storage *Storage) CollectFilesystem(walkPath string) message.FilesystemDirectory {
	filesystemDirectory := message.FilesystemDirectory{
		Directories: map[string]message.FileInfo{},
		Files:       map[string][]message.FileInfo{},
	}

	if walkPath == "" {
		walkPath = "/"
	}

	entries, err := os.ReadDir(walkPath)
	if err != nil {
		return filesystemDirectory
	}

	for _, entry := range entries {
		stat, err := os.Stat(path.Join(walkPath, entry.Name()))
		if err != nil {
			continue
		}

		size := stat.Size() / 1000
		if size == 0 {
			size = 1
		}

		info := message.FileInfo{
			ModTime:  stat.ModTime().Format("02.01.2006 15:04"),
			Size:     size,
			Url:      path.Join(storage.config.Url, storage.config.User.Token),
			Platform: storage.config.Roles.Runner.Platform,
			Hostname: storage.config.Hostname,
		}

		if entry.IsDir() {
			filesystemDirectory.Directories[entry.Name()] = info
		} else {
			filesystemDirectory.Files[entry.Name()] = append(filesystemDirectory.Files[entry.Name()], info)
		}
	}

	return filesystemDirectory
}

func (storage *Storage) GetUrl() string {
	return storage.config.Url + "/" + storage.config.User.Token
}

func (storage *Storage) GetHostnames() map[string]string {
	return map[string]string{
		storage.config.Hostname: path.Join(storage.config.Url, storage.config.User.Token),
	}
}
