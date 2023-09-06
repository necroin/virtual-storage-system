package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
func (server *Server) Start() error {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.instance.Shutdown(ctx)
	}()

	if err := server.instance.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("[App] [Run] [Error] failed start server: %s", err)
	}

	wg.Wait()
	return nil
}

func (server *Server) AddHandler(path string, handler func(http.ResponseWriter, *http.Request)) {
	server.router.HandleFunc(path, handler)
}
