package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"vss/src/logger"
	"vss/src/message"
	"vss/src/roles"
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

func (storage *Storage) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	handlerType := vars["type"]

	data := &message.ClientRequest{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		roles.HandlerFailed(responseWriter, err)
		return
	}

	insertHandler := insertHandlers[handlerType]
	if err := insertHandler(data.Path, data.Name); err != nil {
		roles.HandlerFailed(responseWriter, err)
		return
	}

	roles.HandlerSuccess(responseWriter, "Добавлено")
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("[SelectHandler] failed read path: %s", err)
		return
	}

	selectPath := string(data)
	selectPath = utils.HandleFilesystemPath(selectPath)
	logger.Debug("[SelectHandler] select path: %s", selectPath)

	if err := utils.Compress(selectPath, responseWriter); err != nil {
		logger.Error("[SelectHandler] failed zip data: %s", err)
		return
	}
}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("[DeleteHandler] failed read data: %s", err)
		roles.HandlerFailed(responseWriter, fmt.Errorf("Ошибка чтения данных запроса"))
		return
	}

	deletePath := string(data)
	deletePath = utils.HandleFilesystemPath(deletePath)
	logger.Debug("[DeleteHandler] delete path: %s", deletePath)

	if err := utils.RemoveFile(deletePath); err != nil {
		logger.Error("[DeleteHandler] failed delete: %s", err)
		roles.HandlerFailed(responseWriter, fmt.Errorf("Не удалось выполнить удаление"))
		return
	}
	roles.HandlerSuccess(responseWriter, fmt.Sprintf("%s удалено", path.Base(deletePath)))
}

func (storage *Storage) CopyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	handlerType := vars["type"]

	copyRequest := &message.CopyRequest{}
	copyRequest.NewPath = utils.HandleFilesystemPath(copyRequest.NewPath)
	copyRequest.OldPath = utils.HandleFilesystemPath(copyRequest.OldPath)

	if err := json.NewDecoder(request.Body).Decode(copyRequest); err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[CopyHandler] failed decode response: %s", err)
		return
	}
	logger.Debug("[CopyHandler] request: %#v", copyRequest)

	if handlerType != "dir" {
		copyRequest.NewPath = path.Dir(copyRequest.NewPath)
	}

	response, err := storage.connector.SendPostRequest(copyRequest.SrcUrl+"/storage/select", copyRequest.OldPath)
	if err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[CopyHandler] failed send request: %s", err)
		return
	}
	logger.Debug("[CopyHandler] response: %#v", response)

	if err := utils.Decompress(response.Body, copyRequest.NewPath); err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[CopyHandler] failed unzip file: %s", err)
		return
	}

	roles.HandlerSuccess(responseWriter, "Копирование выполнено")
}

func (storage *Storage) MoveHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	handlerType := vars["type"]

	copyRequest := &message.CopyRequest{}
	copyRequest.NewPath = utils.HandleFilesystemPath(copyRequest.NewPath)
	copyRequest.OldPath = utils.HandleFilesystemPath(copyRequest.OldPath)

	if err := json.NewDecoder(request.Body).Decode(copyRequest); err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed decode response: %s", err)
		return
	}
	logger.Debug("[MoveHandler] request: %#v", copyRequest)

	if handlerType != "dir" {
		copyRequest.NewPath = path.Dir(copyRequest.NewPath)
	}

	selectResponse, err := storage.connector.SendPostRequest(copyRequest.SrcUrl+"/storage/select", copyRequest.OldPath)
	if err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed send request: %s", err)
		return
	}
	logger.Debug("[MoveHandler] select response: %#v", selectResponse)

	if err := utils.Decompress(selectResponse.Body, copyRequest.NewPath); err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed unzip file: %s", err)
		return
	}

	deleteResponse, err := storage.connector.SendPostRequest(copyRequest.SrcUrl+"/storage/delete", copyRequest.OldPath)
	if err != nil {
		roles.HandlerFailed(responseWriter, err)
		logger.Error("[MoveHandler] failed send request: %s", err)
		return
	}
	logger.Debug("[MoveHandler] delete response: %#v", deleteResponse)

	roles.HandlerSuccess(responseWriter, "Перемещение выполнено")
}

func RenameHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data := &message.RenameRequest{}

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		roles.HandlerFailed(responseWriter, err)
		return
	}

	if data.NewName == "" {
		roles.HandlerFailed(responseWriter, fmt.Errorf("Не указано новое имя"))
		return
	}

	if err := utils.Rename(data.Path, data.OldName, data.NewName); err != nil {
		roles.HandlerFailed(responseWriter, err)
		return
	}

	roles.HandlerSuccess(responseWriter, "Переименование выполнено")
}

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	roles.FilesystemHandler(storage, responseWriter, request)
}
