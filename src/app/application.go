package app

import (
	"net/http"
	"sync"
	"vss/src/config"
	"vss/src/roles/router"
	"vss/src/roles/runner"
	"vss/src/roles/storage"
	"vss/src/server"
	"vss/src/settings"
)

type Application struct {
	storageRole *storage.Storage
	runnerRole  *runner.Runner
	routerRole  *router.Router
	config      *config.Config
}

func New() (*Application, error) {
	config, err := config.Load("config.yml")
	if err != nil {
		return nil, err
	}

	return &Application{
		storageRole: nil,
		runnerRole:  nil,
		routerRole:  nil,
		config:      config,
	}, nil
}

func (app *Application) Run() error {
	wg := sync.WaitGroup{}

	server := server.New(app.config.Url)

	server.AddHandler(settings.ServerStatusEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(settings.ServerStatusResponse))
	})

	if app.config.Roles.Storage.Enable {
		storageRole, err := storage.New(app.config.RouterUrl, "storage.db")
		if err != nil {
			return err
		}
		app.storageRole = storageRole
		server.AddHandler(settings.StorageViewEndpoint, storageRole.ViewHandler)
		server.AddHandler(settings.StorageInsertEndpoint, storageRole.InsertHandler)
		server.AddHandler(settings.StorageSelectEndpoint, storageRole.SelectHandler)
		server.AddHandler(settings.StorageUpdateEndpoint, storageRole.UpdateHandler)
		server.AddHandler(settings.StorageDeleteEndpoint, storageRole.DeleteHandler)
		go storageRole.LoadFileSystem()
	}

	if app.config.Roles.Runner.Enable {
		runnerRole, err := runner.New(app.config.RouterUrl)
		if err != nil {
			return err
		}
		app.runnerRole = runnerRole
		server.AddHandler(settings.RouterNotifyEndpoint, runnerRole.NotifyHandler)
		server.AddHandler(settings.RunnerViewEndpoint, runnerRole.ViewHandler)
		server.AddHandler(settings.RunnerOpenEndpoint, runnerRole.OpenFileHandler)
	}

	if app.config.Roles.Router.Enable {
		routerRole, err := router.New()
		if err != nil {
			return err
		}
		app.routerRole = routerRole
		server.AddHandler(settings.RouterViewEndpoint, routerRole.ViewHandler)
		server.AddHandler(settings.RouterTopologyEndpoint, routerRole.GetTopologyHandler)
		server.AddHandler(settings.RouterNotifyEndpoint, routerRole.NotifyHandler)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Start()
	}()

	if err := server.WaitStart(); err != nil {
		return err
	}

	if app.storageRole != nil {
		app.storageRole.NotifyRouter()
	}
	if app.runnerRole != nil {
		app.runnerRole.NotifyRouter()
	}
	if app.routerRole != nil {
		app.routerRole.NotifyRunners()
	}

	wg.Wait()

	return nil
}
