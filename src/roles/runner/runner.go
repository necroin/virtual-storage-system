package runner

import (
	"vss/src/config"
	"vss/src/connector"
	"vss/src/settings"
)

type Runner struct {
	url       string
	routerUrl string
	storages  []string
}

func New(config *config.Config) (*Runner, error) {
	return &Runner{
		url:       config.Url,
		routerUrl: config.RouterUrl,
		storages:  []string{},
	}, nil
}

func (runner *Runner) NotifyRouter() error {
	message := connector.NotifyMessage{
		Type: connector.NotifyMessageRunnerType,
		Url:  runner.url,
	}

	_, err := connector.SendPostRequest(runner.routerUrl+settings.RouterNotifyEndpoint, message)
	return err
}
