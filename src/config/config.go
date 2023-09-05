package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Roles struct {
	Storage StorageRole `yaml:"storage"`
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
	return config, nil
}
