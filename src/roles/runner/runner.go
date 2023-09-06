package runner

type Runner struct {
	routerUrl string
	storages  []string
}

func New(routerUrl string) (*Runner, error) {
	return &Runner{
		routerUrl: routerUrl,
		storages:  []string{},
	}, nil
}

func (runner *Runner) NotifyRouter() {
	// TODO: NotifyRouter
}
