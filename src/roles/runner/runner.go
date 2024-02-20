package runner

import (
	"os"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/settings"
)

type Runner struct {
	url       string
	routerUrl string
	storages  []connector.NotifyMessage
}

func New(config *config.Config) (*Runner, error) {
	return &Runner{
		url:      config.Url,
		storages: []connector.NotifyMessage{},
	}, nil
}

func (runner *Runner) NotifyRouter() error {
	hostname, _ := os.Hostname()

	message := connector.NotifyMessage{
		Type:     connector.NotifyMessageRunnerType,
		Url:      runner.url,
		Hostname: hostname,
	}

	_, err := connector.SendPostRequest(runner.routerUrl+settings.RouterNotifyEndpoint, message)
	return err
}
