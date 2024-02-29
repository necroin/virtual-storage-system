package observer

import (
	"fmt"
	"net"
	"strings"
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
	lanNetwork3 := "0"

	if !strings.Contains(observer.config.Url, "localhost") {
		lanNetworkParts := strings.Split(observer.config.Url, ".")
		lanNetwork3 = lanNetworkParts[2]
	}

	local, err := net.ResolveUDPAddr("udp4", observer.config.Url)
	if err != nil {
		return err
	}

	remote, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("192.168.%s.255:%s", lanNetwork3, observer.config.ListenPort))
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
