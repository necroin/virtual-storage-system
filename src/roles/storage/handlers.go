package storage

import (
	"encoding/json"
	"net/http"
	"vss/src/connector"
	"vss/src/utils/html"

	_ "embed"
)

var (
	//go:embed assets/view.html
	viewTamplate string
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

func (storage *Storage) ViewHandler(responseWriter http.ResponseWriter, r *http.Request) {

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
