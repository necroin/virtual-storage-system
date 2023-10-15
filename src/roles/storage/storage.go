package storage

import (
	"io/ioutil"
	"os"
	"path"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/db"
	"vss/src/settings"
	"vss/src/utils"
)

type Storage struct {
	config   *config.Config
	hostname string
	db       *db.Database
}

func New(config *config.Config, dbPath string) (*Storage, error) {
	db, err := db.New(dbPath)
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	return &Storage{
		config:   config,
		hostname: hostname,
		db:       db,
	}, nil
}

func (storage *Storage) GetRouterToken() (string, error) {
	message := connector.ClientAuth{
		Username: storage.config.User.Username,
		Password: storage.config.User.Password,
	}

	response, err := connector.SendPostRequest(storage.config.RouterUrl+settings.ServerAuthTokenEndpoint, message)
	if err != nil {
		return "", err
	}
	tokenData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	token := string(tokenData)

	return token, nil
}

func (storage *Storage) NotifyRouter() error {
	token, err := storage.GetRouterToken()
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
		storage.config.RouterUrl+utils.FormatTokemizedEndpoint(settings.RouterNotifyEndpoint, token),
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

	entries, _ := os.ReadDir(walkPath)
	for _, entry := range entries {
		stat, _ := os.Stat(path.Join(walkPath, entry.Name()))

		info := connector.FileInfo{
			ModTime: stat.ModTime(),
			Size:    stat.Size(),
			Url:     storage.config.Url,
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
	return storage.config.Url
}

func (storage *Storage) GetMainEndpoint() string {
	return utils.FormatTokemizedEndpoint(settings.StorageMainEndpoint, storage.config.User.Token)
}

func (storage *Storage) GetHostnames() map[string]string {
	return map[string]string{
		storage.hostname: path.Join(storage.config.Url, storage.config.User.Token),
	}
}
