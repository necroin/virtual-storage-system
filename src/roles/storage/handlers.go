package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"vss/src/settings"
	"vss/src/utils/html"

	_ "embed"
)

var (
	//go:embed open.js
	openScript string
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

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	fileSystemMessage := storage.CollectFileSystem(string(msgPath))
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func (storage *Storage) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	walkPath := "/"
	if len(msgPath) != 0 {
		walkPath = string(msgPath)
	}
	walkDirectory := storage.CollectFileSystem(walkPath)

	list := html.NewUnorderedList()

	for _, directory := range walkDirectory.Directories {
		button := html.NewButton(directory, "📁")
		button.SetOnClick(fmt.Sprintf("window.open('%s')", path.Join(walkPath, directory)))
		list.AddElements(button)
	}

	for _, file := range walkDirectory.Files {
		button := html.NewButton(file, "")
		list.AddElements(button)
	}

	head := html.NewHead()
	body := html.NewBody(head)
	body.AddElements(
		html.NewButton("", "←").SetOnClick(fmt.Sprintf("window.open('%s')", path.Join(walkPath, ".."))),
		html.NewText(walkPath),
		list,
		html.NewScript(fmt.Sprintf(openScript, "http://"+storage.url+settings.StorageMainEndpoint)),
	)

	document := html.NewDocument()
	document.AddElements(body)

	responseWriter.Write([]byte(document.ToHTML()))
}
