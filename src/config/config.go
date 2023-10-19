package config

import (
	"fmt"
	"os"
	"time"
	"vss/src/lan"
	"vss/src/settings"
	"vss/src/utils"

	"gopkg.in/yaml.v2"
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

func (log *Log) setDefaults() {
	if log.Path == "" {
		log.Path = "logs/log_" + time.Now().Format("2006-01-02T15:04:05") + ".txt"
	}
}

func (user *User) setDefaults() {
	if user.Username == "" {
		user.Username = settings.DefaultUsername
	}

	if user.Password == "" {
		user.Password = settings.DefaultPassword
	}

	if user.Token == "" {
		user.Token = utils.GenerateSecureToken(10)
	}
}

func (config *Config) setDefaults() {
	if config.Url == "" {
		config.Url = lan.GetMyLanAddr()
	}

	if config.ListenPort == "" {
		config.ListenPort = settings.DefaultListenPort
	}

	config.Log.setDefaults()
	config.User.setDefaults()
}

func Load(path string) (*Config, error) {
	config := &Config{
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
	config.setDefaults()

	configFile, err := os.ReadFile(settings.ConfigPath)
	if err != nil {
		return config, nil
	}

	if err := yaml.Unmarshal(configFile, config); err != nil {
		return nil, fmt.Errorf("[Config] [Error] failed parse config file: %s", err)
	}

	config.setDefaults()

	return config, nil
}
