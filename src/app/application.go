package app

import (
	"fmt"
	"vss/src/config"
)

type Application struct {
	config *config.Config
}

func New() (*Application, error) {
	config, err := config.Load("config.yml")
	if err != nil {
		return nil, err
	}
	return &Application{
		config: config,
	}, nil
}

func (app *Application) Run() {
	fmt.Println(app.config)
}
