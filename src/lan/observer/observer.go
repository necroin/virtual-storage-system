package observer

import (
	"net"
	"time"
	"vss/src/config"
	"vss/src/logger"
	"vss/src/settings"
)

type Observer struct {
	config *config.Config
}

func New(config *config.Config) *Observer {
	return &Observer{
		config: config,
	}
}

func (observer *Observer) Start() error {
	local, err := net.ResolveUDPAddr("udp4", observer.config.Url)
	if err != nil {
		return err
	}

	remote, err := net.ResolveUDPAddr("udp4", "192.168.0.255"+observer.config.ListenPort)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", local, remote)
	if err != nil {
		return err
	}

	go func() {
		for {
			_, err = conn.Write([]byte("observer"))
			if err != nil {
				logger.Debug("%s", err)
			}
			time.Sleep(settings.DefaultLanTimeout)
		}
	}()

	return nil
}
