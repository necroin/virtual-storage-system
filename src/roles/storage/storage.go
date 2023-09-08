package storage

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
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
	filepath.Walk("/", func(path string, info fs.FileInfo, err error) error {
		kind := "file"
		if info.IsDir() {
			kind = "dir"
		}
		err = storage.db.Insert("filesystem", []string{"id", "kind", "path"}, []string{fmt.Sprintf("%d", id), fmt.Sprintf("'%s'", kind), fmt.Sprintf("'%s'", path)})
		if err != nil {
			fmt.Println(err)
		}
		id += 1
		return nil
	})
}

func (storage *Storage) CollectFileSystem() connector.FilesystemDirectory {
	response := storage.db.SelectRequest(&db.Request{
		Table:  "filesystem",
		Fields: []db.Field{{Name: "kind"}, {Name: "path"}},
	})

	fileSystemDir := connector.FilesystemDirectory{
		Directories: map[string]*connector.FilesystemDirectory{},
		Files:       []string{},
	}

	for _, record := range response.Records {
		kind := record.Fields[0].Value
		filePath := record.Fields[1].Value
		path := strings.Split(filePath, string(filepath.Separator))

		workDir := &fileSystemDir
		for i := 0; i < len(path)-1; i++ {
			_, ok := workDir.Directories[path[i]]
			if !ok {
				workDir.Directories[path[i]] = &connector.FilesystemDirectory{
					Directories: map[string]*connector.FilesystemDirectory{},
					Files:       []string{},
				}
			}
			workDir = workDir.Directories[path[i]]
		}
		if kind == "file" {
			workDir.Files = append(workDir.Files, path[len(path)-1])
		}
	}
	return fileSystemDir

}
