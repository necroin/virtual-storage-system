package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/logger"
	"vss/src/settings"
	"vss/src/utils"

	"github.com/gorilla/mux"
)

type Server struct {
	config    *config.Config
	router    *mux.Router
	instance  *http.Server
	connector *connector.Connector
}

func New(config *config.Config, connector *connector.Connector, certificate *tls.Certificate) (*Server, error) {
	router := mux.NewRouter()

	instance := &http.Server{
		Addr:    config.Url,
		Handler: router,
		TLSConfig: &tls.Config{
			ServerName: "vss",
			GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
				logger.Debug("[Server] client requested certificate")
				return certificate, nil
			},
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if len(verifiedChains) > 0 {
					logger.Debug("[Server] Verified certificate chain from peer:")
					for _, certificate := range verifiedChains {
						for i, cert := range certificate {
							logger.Debug("[Server] Cert %d:\n", i)
							logger.Debug(fmt.Sprintf("[Server] %s", utils.CertificateInfo(cert)))
						}
					}
				}
				return nil
			},
		},
	}

	return &Server{
		config:    config,
		router:    router,
		instance:  instance,
		connector: connector,
	}, nil
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
		response, err := server.connector.SendRequest(server.config.Url+settings.ServerStatusEndpoint, []byte(""), http.MethodGet)
		if err != nil {
			logger.Error("[Server] [WaitStart] failed send request: %s", err)
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Error("[Server] [WaitStart] failed read response data: %s", err)
			time.Sleep(settings.ServerWaitStartSleepSeconds * time.Second)
			continue
		}
		if string(data) == settings.ServerStatusResponse {
			return nil
		}
	}
	return fmt.Errorf("[Server] [WaitStart] [Error] failed get server status")
}
