package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"vss/src/connector"
	"vss/src/logger"
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
	vars := mux.Vars(request)
	handlerType := vars["type"]

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
	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("[SelectHandler] failed read path: %s", err)
		return
	}

	if err := utils.Compress(string(data), responseWriter); err != nil {
		logger.Error("[SelectHandler] failed zip data: %s", err)
		return
	}
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
	copyRequest := &connector.CopyRequest{}
	if err := json.NewDecoder(request.Body).Decode(copyRequest); err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[CopyHandler] failed decode response: %s", err)
		return
	}
	logger.Debug("[CopyHandler] request: %#v", copyRequest)

	response, err := connector.SendPostRequest(copyRequest.SrcUrl+"/storage/select", copyRequest.OldPath)
	if err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[CopyHandler] failed send request: %s", err)
		return
	}
	logger.Debug("[CopyHandler] response: %#v", response)

	if err := utils.Decompress(response.Body, copyRequest.NewPath); err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[CopyHandler] failed unzip file: %s", err)
		return
	}

	handlerSuccess(responseWriter, "Копирование выполнено")
}

func (storage *Storage) MoveHandler(responseWriter http.ResponseWriter, request *http.Request) {
	copyRequest := &connector.CopyRequest{}
	if err := json.NewDecoder(request.Body).Decode(copyRequest); err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed decode response: %s", err)
		return
	}
	logger.Debug("[MoveHandler] request: %#v", copyRequest)

	selectResponse, err := connector.SendPostRequest(copyRequest.SrcUrl+"/storage/select", copyRequest.OldPath)
	if err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed send request: %s", err)
		return
	}
	logger.Debug("[MoveHandler] response: %#v", selectResponse)

	if err := utils.Decompress(selectResponse.Body, copyRequest.NewPath); err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed unzip file: %s", err)
		return
	}

	_, err = connector.SendPostRequest(copyRequest.SrcUrl+"/storage/delete", copyRequest.OldPath)
	if err != nil {
		handlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed send request: %s", err)
		return
	}

	handlerSuccess(responseWriter, "Копирование выполнено")
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
