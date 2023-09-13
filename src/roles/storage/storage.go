package storage

import (
	"io/fs"
	"os"
	"path"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/db"
	"vss/src/settings"
)

type Storage struct {
	url       string
	routerUrl string
	db        *db.Database
}

func New(config *config.Config, dbPath string) (*Storage, error) {
	db, err := db.New(dbPath)
	if err != nil {
		return nil, err
	}
	return &Storage{
		url:       config.Url,
		routerUrl: config.RouterUrl,
		db:        db,
	}, nil
}

func (storage *Storage) NotifyRouter() error {
	message := connector.NotifyMessage{
		Type: connector.NotifyMessageStorageType,
		Url:  storage.url,
	}
	_, err := connector.SendPostRequest(storage.routerUrl+settings.RouterNotifyEndpoint, message)
	return err
}

func (storage *Storage) CollectFileSystem(walkPath string) connector.FilesystemDirectory {
	fileSystemDirectory := connector.FilesystemDirectory{
		Directories: map[string]fs.FileInfo{},
		Files:       map[string]fs.FileInfo{},
	}

	if walkPath == "" {
		walkPath = "/"
	}

	entries, _ := os.ReadDir(walkPath)
	for _, entry := range entries {
		stat, _ := os.Stat(path.Join(walkPath, entry.Name()))
		if entry.IsDir() {
			fileSystemDirectory.Directories[entry.Name()] = stat
		} else {
			fileSystemDirectory.Files[entry.Name()] = stat
		}
	}

	return fileSystemDirectory

}
