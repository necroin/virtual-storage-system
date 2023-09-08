package router

import (
	"vss/src/connector"
	"vss/src/settings"
)

type Router struct {
	storages []string
	runners  []string
}

func New() (*Router, error) {
	return &Router{}, nil
}

func (router *Router) NotifyRunner(url string) {
	topology := connector.TopologyMessage{
		Storages: router.storages,
		Runners:  router.runners,
	}
	connector.SendPostRequest(url+settings.RunnerNotifyEndpoint, topology)
}

func (router *Router) NotifyRunners() {
	for _, runner := range router.runners {
		router.NotifyRunner(runner)
	}
}
