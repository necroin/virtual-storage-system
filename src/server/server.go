package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vss/src/settings"

	"github.com/gorilla/mux"
)

type Server struct {
	url      string
	router   *mux.Router
	instance *http.Server
}

func New(url string) *Server {
	router := mux.NewRouter()
	return &Server{
		url:    url,
		router: router,
		instance: &http.Server{
			Addr:    url,
			Handler: router,
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

	server.instance.ListenAndServe()
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request)) {
	server.router.HandleFunc(path, handler)
}

func (server *Server) WaitStart() error {
	client := http.Client{}
	for i := 0; i < settings.ServerWaitStartRepeatCount; i++ {
		response, err := client.Do(&http.Request{
			RequestURI: server.url + settings.ServerStatusEndpoint,
		})
		if err != nil {
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		if string(data) == settings.ServerStatusResponse {
			return nil
		}
	}
	return fmt.Errorf("[Server] [WaitStart] [Error] failed get server status")
}
