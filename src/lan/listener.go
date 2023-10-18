package lan

import (
	"net"
	"vss/src/config"
)

type Listener struct {
	config *config.Config
}

func NewListener(config *config.Config) *Listener {
	return &Listener{
		config: config,
	}
}

func (listener *Listener) Start() {
	netListener, _ := net.Listen("tcp", listener.config.ListenPort)
	go func() {
		for {
			conn, _ := netListener.Accept()
			conn.Close()
		}
	}()

}
