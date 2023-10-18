package config

import (
	"time"
	"vss/src/settings"
	"vss/src/utils"
)

type StorageRole struct {
	Enable bool `yaml:"enable"`
}

type RunnerRole struct {
	Enable bool `yaml:"enable"`
}

type RouterRole struct {
	Enable bool `yaml:"enable"`
}

type Roles struct {
	Storage StorageRole `yaml:"storage"`
	Runner  RunnerRole  `yaml:"runner"`
	Router  RouterRole  `yaml:"router"`
}

type Log struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

type Config struct {
	Url        string `yaml:"url"`
	ListenPort string `yaml:"listen_port"`
	Roles      Roles  `yaml:"roles"`
	Log        Log    `yaml:"log"`
	User       User   `yaml:"user"`
}

func setDefaults(config *Config) {
	if config.Log.Path == "" {
		config.Log.Path = "logs/log_" + time.Now().Format("2006-01-02T15:04:05") + ".txt"
	}

	if config.User.Username == "" {
		config.User.Username = settings.DefaultUsername
	}

	if config.User.Password == "" {
		config.User.Password = settings.DefaultPassword
	}

	if config.User.Token == "" {
		config.User.Token = utils.GenerateSecureToken(10)
	}
}

func Load(path string) (*Config, error) {
	config := &Config{
		Url:        settings.DefaultUrl,
		ListenPort: settings.DefaultListenPort,
		Roles: Roles{
			Storage: StorageRole{
				Enable: true,
			},
			Runner: RunnerRole{
				Enable: true,
			},
			Router: RouterRole{
				Enable: true,
			},
		},
	}

	setDefaults(config)

	return config, nil
}
