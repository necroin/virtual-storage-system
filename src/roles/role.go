package roles

import (
	"encoding/json"
	"io"
	"net/http"
	"vss/src/connector"
	"vss/src/message"
	"vss/src/settings"
)

type Role interface {
	CollectFilesystem(walkPath string) message.FilesystemDirectory
	GetUrl() string
	GetHostnames() map[string]string
}

func FilesystemHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := io.ReadAll(request.Body)

	fileSystemMessage := role.CollectFilesystem(string(msgPath))
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func DevicesHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	hostnames := role.GetHostnames()
	json.NewEncoder(responseWriter).Encode(hostnames)
}

func GetRouterToken(connector *connector.Connector, url string, username string, password string) (string, error) {
	message := message.ClientAuthRequest{
		Username: username,
		Password: password,
	}

	response, err := connector.SendPostRequest(url+settings.ServerAuthTokenEndpoint, message)
	if err != nil {
		return "", err
	}
	tokenData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	token := string(tokenData)

	return token, nil
}

func HandlerFailed(responseWriter http.ResponseWriter, err error) {
	json.NewEncoder(responseWriter).Encode(message.StatusBarResponse{
		Status: settings.ExplorerStatusBarFail,
		Text:   err.Error(),
	})
}

func HandlerSuccess(responseWriter http.ResponseWriter, text string) {
	json.NewEncoder(responseWriter).Encode(message.StatusBarResponse{
		Status: settings.ExplorerStatusBarSuccess,
		Text:   text,
	})
}
