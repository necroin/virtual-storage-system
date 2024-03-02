package runner

import (
	"os"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"
)

type Runner struct {
	config   *config.Config
	hostname string
	storages []connector.NotifyMessage
}

func New(config *config.Config) (*Runner, error) {
	hostname, _ := os.Hostname()

	return &Runner{
		config:   config,
		hostname: hostname,
		storages: []connector.NotifyMessage{},
	}, nil
}

func (runner *Runner) NotifyRouter(url string) error {
	token, err := roles.GetRouterToken(url, runner.config.User.Username, runner.config.User.Password)
	if err != nil {
		return err
	}

	message := connector.NotifyMessage{
		Type:     connector.NotifyMessageRunnerType,
		Url:      runner.config.Url,
		Hostname: runner.hostname,
		Token:    runner.config.User.Token,
		Platform: runner.config.Roles.Runner.Platform,
	}

	_, err = connector.SendPostRequest(
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
