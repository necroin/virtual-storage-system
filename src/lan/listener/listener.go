package listener

import (
	"fmt"
	"net"
	"vss/src/config"
	"vss/src/logger"
)

type Listener struct {
	config *config.Config
}

func New(config *config.Config) *Listener {
	return &Listener{
		config: config,
	}
}

func (listener *Listener) Start() (chan string, error) {
	result := make(chan string)

	udpAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%s", listener.config.ListenPort))
	if err != nil {
		return nil, err
	}

	netListener, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		buf := make([]byte, 1024)
		for {
			_, addr, err := netListener.ReadFromUDP(buf)
			if err != nil {
				panic(err)
			}
			logger.Debug("[Lan Listener] %s", addr.String())
			result <- addr.String()
		}
	}()
	return result, nil
}
