package settings

import "time"

const (
	ServerWaitStartRepeatCount  = 10
	ServerWaitStartSleepSeconds = 1
	ServerStatusResponse        = "OK"
)

const (
	ConfigPath = "config.yml"
)

const (
	DefaultPort                  = ":3301"
	DefaultListenPort            = "3311"
	DefaultUsername              = "admin"
	DefaultPassword              = "admin"
	DefaultInstanceRemoveSeconds = 30
)

const (
	DefaultLanTimeout = 10 * time.Second
)
