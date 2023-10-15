package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vss/src/connector"
	"vss/src/settings"
	"vss/src/utils/html"

	"github.com/gorilla/mux"
)

func (server *Server) StatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(settings.ServerStatusResponse))
}

func (server *Server) AuthHandler(responseWriter http.ResponseWriter, request *http.Request) {
	style := html.NewTag("style").AddElements(html.NewText(settings.GetAuthenticationStyle())).AddAttribute(html.NewAttribute("type", "text/css"))
	script := html.NewScript(fmt.Sprintf(settings.GetAuthenticationScript(), server.url))

	result := fmt.Sprintf(
		settings.GetAuthenticationTemlate(),
		style.ToHTML(),
		script.ToHTML(),
	)

	responseWriter.Write([]byte(result))
}

func (server *Server) AuthTokenHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data := &connector.ClientAuth{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}

	if data.Username == server.config.User.Username && data.Password == server.config.User.Password {
		responseWriter.Write([]byte(server.config.User.Token))
	}
}

func (server *Server) TokenizedHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		token, ok := vars["token"]
		if !ok {
			return
		}

		if token == server.config.User.Token {
			handler(responseWriter, request)
		}
	}
}
