package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

func (storage *Storage) ViewHandler(responseWriter http.ResponseWriter, r *http.Request) {
	fileSystemView := storage.CollectFileSystem()

	files := fmt.Sprintf("<head>Files</head>\n<body>\n<ul>\n<li>%s</li>\n</ul>\n</body>\n", strings.Join(fileSystemView.Files, "</li>\n<li>"))
	dirs := []string{}
	for name := range fileSystemView.Directories[""].Directories {
		fmt.Println(name)
		dirs = append(dirs, name)
	}
	directories := fmt.Sprintf("<head>Directories</head>\n<body>\n<ul>\n<li>%s</li>\n</ul>\n</body>", strings.Join(dirs, "</li>\n<li>"))
	fmt.Println(fmt.Sprintf(viewTamplate, files+directories))
	responseWriter.Write([]byte(fmt.Sprintf(viewTamplate, files+directories)))
}
