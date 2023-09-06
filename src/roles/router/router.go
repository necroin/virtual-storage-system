package router

type Router struct {
	storages []string
	runners  []string
}

func New() (*Router, error) {
	return &Router{}, nil
}

func (router *Router) NotifyRunners() {
	// TODO: NotifyRunners
}
