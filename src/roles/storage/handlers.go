package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"vss/src/buffer"
	"vss/src/connector"
	"vss/src/roles"
	"vss/src/utils"

	_ "embed"

	"github.com/gorilla/mux"
)

var (
	insertHandlers = map[string]func(string, string) error{
		"dir": func(path string, name string) error {
			if name == "" {
				name = "Новая папка"
			}
			utils.CreateNewDirectory(path, name)
			return nil
		},
		"file": func(path string, name string) error {
			if name == "" {
				return fmt.Errorf("Имя файла не указано")
			}
			utils.CreateNewFile(path, name)
			return nil

		},
		"textFile": func(path string, name string) error {
			if name == "" {
				name = "Текстовый документ"
			}
			name += ".txt"
			utils.CreateNewFile(path, name)
			return nil
		},
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

	data := &connector.ClientRequest{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}

	insertHandler := insertHandlers[handlerType]
	if err := insertHandler(data.Path, data.Name); err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}

	responseWriter.Write([]byte("Добавлено"))
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data, _ := ioutil.ReadAll(request.Body)
	deletePath := string(data)
	err := utils.RemoveFile(deletePath)
	if err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}
	responseWriter.Write([]byte(fmt.Sprintf("%s удалено", path.Base(deletePath))))
}

func (storage *Storage) CopyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	handlerType := msgVars["type"]

	data, _ := ioutil.ReadAll(request.Body)
	copyPath := string(data)

	copyHandler := copyHandlers[handlerType]
	copyHandler(copyPath, handlerType)
	responseWriter.Write([]byte(fmt.Sprintf("%s добавлен в буффер копирования", path.Base(copyPath))))
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
