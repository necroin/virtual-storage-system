package tests

import (
	"fmt"
	"testing"
	"time"
	"vss/src/config"
	"vss/src/lan/listener"
	"vss/src/lan/observer"
)

func TestLan(t *testing.T) {
	config, _ := config.LoadConfig()
	listener := listener.New(config)
	addrs, err := listener.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	observer := observer.New(config)
	observer.Start()
	go func() {
		for {
			fmt.Printf("[TEST] %s\n", <-addrs)
		}
	}()
	time.Sleep(time.Second * 10)
}
