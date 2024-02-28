package storage

import (
	"os"
	"path"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"
)

type Storage struct {
	config   *config.Config
	hostname string
}

func New(config *config.Config) (*Storage, error) {
	hostname, _ := os.Hostname()

	return &Storage{
		config:   config,
		hostname: hostname,
	}, nil
}

func (storage *Storage) NotifyRouter(url string) error {
	token, err := roles.GetRouterToken(url, storage.config.User.Username, storage.config.User.Password)
	if err != nil {
		return err
	}

	message := connector.NotifyMessage{
		Type:     connector.NotifyMessageStorageType,
		Url:      storage.config.Url,
		Hostname: storage.hostname,
		Token:    storage.config.User.Token,
	}
	_, err = connector.SendPostRequest(
		url+utils.FormatTokemizedEndpoint(settings.RouterNotifyEndpoint, token),
		message,
	)
	return err
}

func (storage *Storage) CollectFileSystem(walkPath string) connector.FilesystemDirectory {
	fileSystemDirectory := connector.FilesystemDirectory{
		Directories: map[string]connector.FileInfo{},
		Files:       map[string]connector.FileInfo{},
	}

	if walkPath == "" {
		walkPath = "/"
	}

	entries, err := os.ReadDir(walkPath)
	if err != nil {
		return fileSystemDirectory
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

		info := connector.FileInfo{
			ModTime:  stat.ModTime().Format("02.01.2006 15:04"),
			Size:     size,
			Url:      path.Join(storage.config.Url, storage.config.User.Token),
			Platform: storage.config.Roles.Runner.Platform,
		}

		if entry.IsDir() {
			fileSystemDirectory.Directories[entry.Name()] = info
		} else {
			fileSystemDirectory.Files[entry.Name()] = info
		}
	}

	return fileSystemDirectory
}

func (storage *Storage) GetUrl() string {
	return storage.config.Url + "/" + storage.config.User.Token
}

func (storage *Storage) GetHostnames() map[string]string {
	return map[string]string{
		storage.hostname: path.Join(storage.config.Url, storage.config.User.Token),
	}
}
