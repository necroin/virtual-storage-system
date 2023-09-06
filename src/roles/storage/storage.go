package storage

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"vss/src/db"
)

type Storage struct {
	routerUrl string
	db        *db.Database
}

func New(routerUrl string, dbPath string) (*Storage, error) {
	db, err := db.New(dbPath)
	if err != nil {
		return nil, err
	}
	return &Storage{
		routerUrl: routerUrl,
		db:        db,
	}, nil
}

func (storage *Storage) NotifyRouter() {
	// TODO: NotifyRouter
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
