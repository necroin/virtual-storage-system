package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"vss/src/connector"
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

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, r *http.Request) {
	fileSystemMessage := storage.CollectFileSystem()
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func buildFileSystemHtml(name string, exploreDir *connector.FilesystemDirectory) string {
	files := html.NewBody(html.NewHead())
	filesList := html.NewUnorderedList()
	for _, file := range exploreDir.Files {
		filesList.Add(html.NewText(file))
	}
	files.Add(filesList)

	directories := html.NewBody(html.NewHead())
	directoriesList := html.NewUnorderedList()
	for dirName, dir := range exploreDir.Directories {
		dirHTML := buildFileSystemHtml(dirName, dir)
		directoriesList.Add(html.NewText(dirHTML))
	}
	directories.Add(directoriesList)

	return html.NewBody(html.NewHead().Add(html.NewText(name))).Add(files, directories).ToHTML()
}

func (storage *Storage) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {

	fileSystemView := storage.CollectFileSystem()

	document := html.NewDocument()
	document.Add(
		html.NewBody(
			html.NewHead().Add(html.NewText("FileSystem")),
		).Add(
			html.NewText(buildFileSystemHtml("/", &fileSystemView)),
		),
	)

	responseWriter.Write([]byte(document.ToHTML()))
}

func (storage *Storage) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)
	walkPath := "/"
	if len(msgPath) != 0 {
		walkPath = string(msgPath)
	}

	list := html.NewUnorderedList()
	entries, _ := os.ReadDir(walkPath)
	for _, e := range entries {
		list.Add(
			html.NewButton(e.Name()).SetOnClick(fmt.Sprintf("open(/%s/)", e.Name())),
		)
	}
	document := html.NewDocument()
	body := html.NewBody(html.NewHead()).Add(list).Add(html.NewScript(fmt.Sprintf(openScript, storage.url+settings.StorageMainEndpoint)))
	document.Add(body)
	responseWriter.Write([]byte(document.ToHTML()))
}
