package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/settings"

	"github.com/gorilla/mux"
)

type Server struct {
	config   *config.Config
	url      string
	router   *mux.Router
	instance *http.Server
}

func New(config *config.Config) *Server {
	router := mux.NewRouter()

	return &Server{
		config: config,
		url:    config.Url,
		router: router,
		instance: &http.Server{
			Addr:    config.Url,
			Handler: router,
			TLSConfig: &tls.Config{
				ServerName: "localhost",
			},
		},
	}
}
func (server *Server) Start() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.instance.Shutdown(ctx)
	}()

	server.instance.ListenAndServeTLS("certificates/vss.crt", "certificates/vss.key")
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	server.router.HandleFunc(path, handler).Methods(methods...)
}

func (server *Server) WaitStart() error {
	for i := 0; i < settings.ServerWaitStartRepeatCount; i++ {
		response, err := connector.SendRequest(server.url+settings.ServerStatusEndpoint, []byte(""), http.MethodGet)
		if err != nil {
			fmt.Println(err)
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		if string(data) == settings.ServerStatusResponse {
			return nil
		}
	}
	return fmt.Errorf("[Server] [WaitStart] [Error] failed get server status")
}
