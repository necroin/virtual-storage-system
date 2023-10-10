package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"vss/src/buffer"
	"vss/src/roles"
	"vss/src/utils"

	_ "embed"

	"github.com/gorilla/mux"
)

var (
	insertHandlers = map[string]func(string){
		"dir":      func(path string) { utils.CreateNewDirectory(path, "Новая папка") },
		"textFile": func(path string) { utils.CreateNewFile(path, "Текстовый документ.txt") },
	}
	copyHandlers = map[string]func(string, string){
		"file": func(filePath string, fileType string) { buffer.SetFile(filePath, fileType) },
		"text": func(text string, _ string) {},
	}
	pasteHandlers = map[string]func(string, string) error{
		"dir":  func(srcPath string, dstPath string) error { return nil },
		"file": func(srcPath string, dstPath string) error { return utils.CopyFile(srcPath, dstPath) },
	}
)

func (storage *Storage) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	handlerType := msgVars["type"]

	path, _ := ioutil.ReadAll(request.Body)

	insertHandler := insertHandlers[handlerType]
	insertHandler(string(path))
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)
	utils.RemoveFile(string(msgPath))
}

func (storage *Storage) CopyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	handlerType := msgVars["type"]

	path, _ := ioutil.ReadAll(request.Body)

	copyHandler := copyHandlers[handlerType]
	copyHandler(string(path), handlerType)
	responseWriter.Write([]byte(fmt.Sprintf("%s добавлен в буффер копирования", path)))
}

func (storage *Storage) PasteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	dstPath, _ := ioutil.ReadAll(request.Body)
	srcPath, handlerType, err := buffer.GetFile()
	if err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}

	pasteHandler := pasteHandlers[handlerType]
	err = pasteHandler(srcPath, path.Join(string(dstPath), path.Base(srcPath)))
	if err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}
	responseWriter.Write([]byte("Вставка выполнена"))
}

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.FilesystemHandler(storage, responseWriter, request)
}

func (storage *Storage) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.MainHandler(storage, responseWriter, request)
}
