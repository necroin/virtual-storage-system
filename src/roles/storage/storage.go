package storage

import (
	"os"
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
		Directories: []string{},
		Files:       []string{},
	}

	if walkPath == "" {
		walkPath = "/"
	}

	entries, _ := os.ReadDir(walkPath)
	for _, entry := range entries {
		if entry.IsDir() {
			fileSystemDirectory.Directories = append(fileSystemDirectory.Directories, entry.Name())
		} else {
			fileSystemDirectory.Files = append(fileSystemDirectory.Files, entry.Name())
		}
	}

	return fileSystemDirectory

}
