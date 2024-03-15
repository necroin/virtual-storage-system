package runner

import (
	"os"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/message"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"
)

type Runner struct {
	config    *config.Config
	connector *connector.Connector
	hostname  string
	storages  []message.NotifyMessage
}

func New(config *config.Config, connector *connector.Connector) (*Runner, error) {
	hostname, _ := os.Hostname()

	return &Runner{
		config:    config,
		connector: connector,
		hostname:  hostname,
		storages:  []message.NotifyMessage{},
	}, nil
}

func (runner *Runner) NotifyRouter(url string) error {
	token, err := roles.GetRouterToken(runner.connector, url, runner.config.User.Username, runner.config.User.Password)
	if err != nil {
		return err
	}

	message := message.NotifyMessage{
		Type:      message.NotifyMessageRunnerType,
		Url:       runner.config.Url,
		Hostname:  runner.hostname,
		Token:     runner.config.User.Token,
		Platform:  runner.config.Roles.Runner.Platform,
		Timestamp: time.Now().UnixNano(),
	}

	_, err = runner.connector.SendPostRequest(
		url+utils.FormatTokemizedEndpoint(settings.RouterNotifyEndpoint, token),
		message,
	)
	return err
}

func (runner *Runner) GetRunCommand(path string) (string, []string) {
	if runner.config.Roles.Runner.Platform == "windows" {
		return "cmd", []string{"/C", path}
	}
	return "open", []string{path}
}
