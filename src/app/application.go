package app

import (
	"fmt"
	"sync"
	"vss/src/config"
	"vss/src/roles/router"
	"vss/src/roles/runner"
	"vss/src/roles/storage"
	"vss/src/server"
	"vss/src/settings"

	"github.com/necroin/golibs/metrics"
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
	metricsRegistry := metrics.NewRegistry()

	server := server.New(app.config)

	server.AddHandler(settings.ServerStatusEndpoint, server.StatusHandler, "GET")
	server.AddHandler(settings.ServerAuthEndpoint, server.AuthHandler, "GET")
	server.AddHandler(settings.ServerAuthTokenEndpoint, server.AuthTokenHandler, "POST")

	if app.config.Roles.Storage.Enable {
		storageRole, err := storage.New(app.config, "storage.db")
		if err != nil {
			return err
		}
		app.storageRole = storageRole

		server.AddHandler(settings.StorageMainEndpoint, storageRole.MainHandler, "POST", "GET")
		server.AddHandler(settings.StorageFilesystemEndpoint, storageRole.FilesystemHandler, "GET")
		server.AddHandler(settings.StorageInsertEndpoint, storageRole.InsertHandler, "POST")
		server.AddHandler(settings.StorageSelectEndpoint, storageRole.SelectHandler, "POST")
		server.AddHandler(settings.StorageUpdateEndpoint, storageRole.UpdateHandler, "POST")
		server.AddHandler(settings.StorageDeleteEndpoint, storageRole.DeleteHandler, "POST")
		server.AddHandler(settings.StorageCopyEndpoint, storageRole.CopyHandler, "POST")
		server.AddHandler(settings.StoragePasteEndpoint, storageRole.PasteHandler, "POST")
	}

	// if app.config.Roles.Runner.Enable {
	// 	runnerRole, err := runner.New(app.config)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	app.runnerRole = runnerRole

	// 	server.AddHandler(settings.RunnerTopologyEndpoint, runnerRole.GetTopologyHandler, "GET")
	// 	server.AddHandler(settings.RunnerNotifyEndpoint, runnerRole.NotifyHandler, "POST")
	// 	server.AddHandler(settings.RunnerOpenEndpoint, runnerRole.OpenFileHandler, "POST")
	// }

	if app.config.Roles.Router.Enable {
		routerRole, err := router.New(app.config)
		if err != nil {
			return err
		}
		app.routerRole = routerRole

		server.AddHandler(settings.RouterMainEndpoint, server.TokenizedHandler(routerRole.MainHandler), "POST", "GET")
		server.AddHandler(settings.RouterTopologyEndpoint, server.TokenizedHandler(routerRole.GetTopologyHandler), "GET")
		server.AddHandler(settings.RouterNotifyEndpoint, server.TokenizedHandler(routerRole.NotifyHandler), "POST")
		server.AddHandler(settings.RouterInsertEndpoint, server.TokenizedHandler(routerRole.InsertHandler), "POST")
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Start()
	}()

	if err := server.WaitStart(); err != nil {
		return err
	}
	fmt.Println("Server started.")

	if app.storageRole != nil {
		if err := app.storageRole.NotifyRouter(); err != nil {
			return err
		}
	}
	if app.runnerRole != nil {
		if err := app.runnerRole.NotifyRouter(); err != nil {
			return err
		}
	}

	server.AddHandler(settings.ServerMetricsEndpoint, metricsRegistry.Handler().ServeHTTP, "GET")

	wg.Wait()

	return nil
}
