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
	"vss/src/settings"
	"vss/src/utils"

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
			return utils.CreateNewFile(path, name)

		},
		"textFile": func(path string, name string) error {
			if name == "" {
				name = "Текстовый документ"
			}
			name += ".txt"
			return utils.CreateNewFile(path, name)
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

func handlerFailed(responseWriter http.ResponseWriter, err error) {
	json.NewEncoder(responseWriter).Encode(connector.StatusBarResponse{
		Status: settings.ExplorerStatusBarFail,
		Text:   err.Error(),
	})
}

func handlerSuccess(responseWriter http.ResponseWriter, text string) {
	json.NewEncoder(responseWriter).Encode(connector.StatusBarResponse{
		Status: settings.ExplorerStatusBarSuccess,
		Text:   text,
	})
}

func (storage *Storage) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	handlerType := msgVars["type"]

	data := &connector.ClientRequest{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		handlerFailed(responseWriter, err)
		return
	}

	insertHandler := insertHandlers[handlerType]
	if err := insertHandler(data.Path, data.Name); err != nil {
		handlerFailed(responseWriter, err)
		return
	}

	handlerSuccess(responseWriter, "Добавлено")
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handlerFailed(responseWriter, fmt.Errorf("Ошибка чтения данных запроса"))
		return
	}

	deletePath := string(data)
	if err := utils.RemoveFile(deletePath); err != nil {
		handlerFailed(responseWriter, err)
		return
	}
	handlerSuccess(responseWriter, fmt.Sprintf("%s удалено", path.Base(deletePath)))
}

func (storage *Storage) CopyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	handlerType := msgVars["type"]

	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handlerFailed(responseWriter, fmt.Errorf("Ошибка чтения данных запроса"))
		return
	}
	copyPath := string(data)

	copyHandler := copyHandlers[handlerType]
	copyHandler(copyPath, handlerType)
	handlerSuccess(responseWriter, fmt.Sprintf("%s добавлен в буффер копирования", path.Base(copyPath)))
}

func (storage *Storage) PasteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	dstPath, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handlerFailed(responseWriter, fmt.Errorf("Ошибка чтения данных запроса"))
		return
	}
	srcPath, handlerType, err := buffer.GetFile()
	if err != nil {
		handlerFailed(responseWriter, err)
		return
	}

	pasteHandler := pasteHandlers[handlerType]
	if err := pasteHandler(srcPath, path.Join(string(dstPath), path.Base(srcPath))); err != nil {
		handlerFailed(responseWriter, err)
		return
	}
	handlerSuccess(responseWriter, "Вставка выполнена")
}

func RenameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data := &connector.RenameRequest{}

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		handlerFailed(responseWriter, err)
		return
	}

	if data.NewName == "" {
		handlerFailed(responseWriter, fmt.Errorf("Не указано новое имя"))
		return
	}

	if err := utils.Rename(data.Path, data.OldName, data.NewName); err != nil {
		handlerFailed(responseWriter, err)
		return
	}

	handlerSuccess(responseWriter, "Переименование выполнено")
}

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.FilesystemHandler(storage, responseWriter, request)
}
