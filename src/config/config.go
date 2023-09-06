package config

import (
	"fmt"
	"os"

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

type Config struct {
	Url       string `yaml:"url"`
	RouterUrl string `yaml:"router"`
	Roles     Roles  `yaml:"roles"`
}

func Load(path string) (*Config, error) {
	configFileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[Config] [Load] [Error] failed read config file: %s", err)
	}
	config := &Config{}
	if err := yaml.Unmarshal(configFileContent, config); err != nil {
		return nil, fmt.Errorf("[Config] [Load] [Error] failed parse config file: %s", err)
	}

	if config.RouterUrl == "" {
		config.RouterUrl = config.Url
	}

	return config, nil
}
