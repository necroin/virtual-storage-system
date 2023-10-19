package observer

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"vss/src/config"
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

func (observer *Observer) Start() chan string {
	result := make(chan string)
	for i := 0; i < 256; i++ {
		go func(addr string, port string) {
			ip := addr + port
			for {
				conn, err := net.DialTimeout("tcp", ip, settings.DefaultLanTimeout)
				if err == nil {
					fmt.Println(addr)
					result <- addr
					conn.Close()
				}
				time.Sleep(settings.DefaultLanTimeout)
			}

		}(settings.DefaultLanIP+strconv.Itoa(i), settings.DefaultListenPort)
	}
	return result
}

func (observer *Observer) GetRouters() {
}
