package storage

import (
	"net/http"
	"vss/src/db"
)

func (storage *Storage) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) ViewHandler(responseWriter http.ResponseWriter, r *http.Request) {
	response := storage.db.SelectRequest(&db.Request{
		Table:  "filesystem",
		Fields: []db.Field{{Name: "path"}},
	})

	for _, record := range response.Records {
		responseWriter.Write([]byte(record.Fields[0].Value + "\n"))
	}
}
