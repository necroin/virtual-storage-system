package runner

import (
	"fmt"
	"os"
	"sync"
	"time"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/logger"
	"vss/src/message"
	"vss/src/roles"
	"vss/src/settings"
	"vss/src/utils"

	"github.com/necroin/golibs/utils/winapi"
	"github.com/necroin/golibs/utils/winutils"
	"github.com/necroin/golibs/winappstream"
)

type StreamSession struct {
	app            *winappstream.App
	handler        winappstream.HttpImageCaptureHandler
	lastHandleTime time.Time
}

type Runner struct {
	config         *config.Config
	connector      *connector.Connector
	hostname       string
	streamSessions map[int]*StreamSession
	sessuinMutex   sync.Mutex
}

func New(config *config.Config, connector *connector.Connector) (*Runner, error) {
	hostname, _ := os.Hostname()

	return &Runner{
		config:         config,
		connector:      connector,
		hostname:       hostname,
		streamSessions: map[int]*StreamSession{},
		sessuinMutex:   sync.Mutex{},
	}, nil
}

func (runner *Runner) NotifyRouter(url string) error {
	token, err := roles.GetRouterToken(runner.connector, url, runner.config.User.Username, runner.config.User.Password)
	if err != nil {
		return err
	}

	message := message.Notify{
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

func (runner *Runner) GetProcessPidByParentId(pid int) (int, error) {
	processes, err := winutils.GetAllProcesses()
	if err != nil {
		return 0, fmt.Errorf("[Runner] [OpenFileHandler] failed get all processes")
	}

	childProcesses := winutils.FindProcessesByParentPid(processes, winapi.ProcessId(pid))
	if len(childProcesses) == 0 {
		return 0, fmt.Errorf("[Runner] [OpenFileHandler] failed find child processes for pid %d", pid)
	}
	logger.Debug("[Runner] [OpenFileHandler] childs of %d pid: %s", pid, childProcesses)
	appPid := childProcesses[0].Pid
	logger.Debug("[Runner] [OpenFileHandler] app process started with pid: %d", appPid)
	return int(appPid), nil
}
