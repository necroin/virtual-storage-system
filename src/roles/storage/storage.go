package storage

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"vss/src/db"
)

type Storage struct {
	db *db.Database
}

func New(path string) (*Storage, error) {
	db, err := db.New(path)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (storage *Storage) InsertHandler(w http.ResponseWriter, r *http.Request) {
	storage.db.InsertHandler(r.Body, w)
}

func (storage *Storage) SelectHandler(w http.ResponseWriter, r *http.Request) {
	storage.db.InsertHandler(r.Body, w)
}

func (storage *Storage) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	storage.db.InsertHandler(r.Body, w)
}

func (storage *Storage) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	storage.db.InsertHandler(r.Body, w)
}

func (storage *Storage) ViewHandler(responseWriter http.ResponseWriter, r *http.Request) {
	response := storage.db.Select(&db.Request{
		Table:  "filesystem",
		Fields: []db.Field{{Name: "path"}},
	})

	for _, record := range response.Records {
		responseWriter.Write([]byte(record.Fields[0].Value + "\n"))
	}
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
