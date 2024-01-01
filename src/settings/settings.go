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
	DefaultPort       = ":3301"
	DefaultListenPort = ":3311"
	DefaultUsername   = "admin"
	DefaultPassword   = "admin"
)

const (
	DefaultLanIP      = "192.168.0."
	DefaultLanTimeout = 10 * time.Second
)
