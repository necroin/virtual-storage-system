package main

import (
	"vss/src/app"
	"vss/src/config"
	"vss/src/logger"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := logger.Configure(config.Log.Enable, config.Log.Path, config.Log.Level); err != nil {
		panic(err)
	}

	app, err := app.New(config)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
