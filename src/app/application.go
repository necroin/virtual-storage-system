package app

import (
	"vss/src/config"
	"vss/src/roles/storage"
	"vss/src/server"
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

func (app *Application) Run() error {
	server := server.New(app.config.Url)

	if app.config.Roles.Storage.Enable {
		storageRole, err := storage.New("storage.db")
		if err != nil {
			return err
		}
		server.AddHandler("/view", storageRole.ViewHandler)
		server.AddHandler("/insert", storageRole.InsertHandler)
		server.AddHandler("/select", storageRole.SelectHandler)
		server.AddHandler("/update", storageRole.UpdateHandler)
		server.AddHandler("/delete", storageRole.DeleteHandler)
		go storageRole.LoadFileSystem()
	}

	if err := server.Start(); err != nil {
		return err
	}

	return nil
}
