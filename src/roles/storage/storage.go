package storage

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
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

func (storage *Storage) LoadFileSystem() {
	id := int64(0)
	filepath.Walk("/", func(filePath string, info fs.FileInfo, err error) error {
		err = storage.db.Insert("filesystem", []string{"id", "path"}, []string{fmt.Sprintf("%d", id), fmt.Sprintf("'%s'", path.Join(filePath, info.Name()))})
		if err != nil {
			fmt.Println(err)
		}
		id += 1
		return nil
	})
}
