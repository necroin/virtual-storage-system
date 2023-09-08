package config

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
	config := &Config{
		Url: "localhost:3301",
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

	if config.RouterUrl == "" {
		config.RouterUrl = config.Url
	}

	return config, nil
}
