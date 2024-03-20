package server

import (
	"encoding/json"
	"net/http"
	"text/template"
	"vss/src/logger"
	"vss/src/message"
	"vss/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) StatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(settings.ServerStatusResponse))
}

func (server *Server) PageHandler(responseWriter http.ResponseWriter, htmlPage string, pageInfo message.PageInfo) {
	pageInfo.Url = server.config.Url
	pageInfo.BarHomeIcon = settings.BarHomeIcon
	pageInfo.BarFilesystemIcon = settings.BarFilesystemIcon
	pageInfo.BarSettingsIcon = settings.BarSettingsIcon
	pageTemplate, _ := template.New("HtmlPage").Parse(htmlPage)
	pageTemplate.Execute(responseWriter, pageInfo)
}

func (server *Server) AuthHandler(responseWriter http.ResponseWriter, request *http.Request) {
	server.PageHandler(
		responseWriter,
		settings.GetAuthenticationTemlate(),
		message.PageInfo{
			Style:  settings.GetAuthenticationStyle(),
			Script: settings.GetAuthenticationScript(),
		},
	)
}

func (server *Server) AuthTokenHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data := &message.ClientAuth{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		logger.Error("[Server] [AuthTokenHandler] failed decode message: %s", err)
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
			responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
			handler(responseWriter, request)
		}

	}
}

func (server *Server) HomeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	server.PageHandler(
		responseWriter,
		settings.GetHomePage(),
		message.PageInfo{
			Style:  settings.GetHomePageStyle(),
			Script: settings.GetHomePageScript(),
			Token:  server.config.User.Token,
		},
	)
}

func (server *Server) SettingsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	server.PageHandler(
		responseWriter,
		settings.GetSettingsPage(),
		message.PageInfo{
			Style:  settings.GetSettingsPageStyle(),
			Script: settings.GetSettingsPageScript(),
			Token:  server.config.User.Token,
		},
	)
}
