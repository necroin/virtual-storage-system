package server

import (
	"encoding/json"
	"net/http"
	"text/template"
	"vss/src/connector"
	"vss/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) StatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(settings.ServerStatusResponse))
}

func (server *Server) PageHandler(responseWriter http.ResponseWriter, htmlPage string, pageInfo connector.PageInfo) {
	pageInfo.Url = server.url
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
		connector.PageInfo{
			Style:  settings.GetAuthenticationStyle(),
			Script: settings.GetAuthenticationScript(),
		},
	)
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
			responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
			handler(responseWriter, request)
		}

	}
}

func (server *Server) HomeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	server.PageHandler(
		responseWriter,
		settings.GetHomeTemplate(),
		connector.PageInfo{
			Style:  settings.GetHomeStyle(),
			Script: settings.GetHomeScript(),
			Token:  server.config.User.Token,
		},
	)
}
