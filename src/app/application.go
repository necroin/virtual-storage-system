package app

import (
	"fmt"
	"sync"
	"vss/src/config"
	"vss/src/lan/listener"
	"vss/src/lan/observer"
	"vss/src/logger"
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
	server      *server.Server
	lanListener *listener.Listener
	lanObserver *observer.Observer
}

func New() (*Application, error) {
	config, err := config.Load(settings.ConfigPath)
	if err != nil {
		return nil, err
	}

	if err := logger.Configure(config.Log.Enable, config.Log.Path, config.Log.Level); err != nil {
		return nil, err
	}

	server := server.New(config)

	server.AddHandler(settings.ServerStatusEndpoint, server.StatusHandler, "GET")
	server.AddHandler(settings.ServerAuthEndpoint, server.AuthHandler, "GET")
	server.AddHandler(settings.ServerAuthTokenEndpoint, server.AuthTokenHandler, "POST")

	var storageRole *storage.Storage = nil
	if config.Roles.Storage.Enable {
		storageRole, err = storage.New(config)
		if err != nil {
			return nil, err
		}

		server.AddHandler(settings.StorageFilesystemEndpoint, storageRole.FilesystemHandler, "POST", "GET")
		server.AddHandler(settings.StorageInsertEndpoint, storageRole.InsertHandler, "POST")
		server.AddHandler(settings.StorageSelectEndpoint, storageRole.SelectHandler, "POST")
		server.AddHandler(settings.StorageUpdateEndpoint, storageRole.UpdateHandler, "POST")
		server.AddHandler(settings.StorageDeleteEndpoint, storageRole.DeleteHandler, "POST")
		server.AddHandler(settings.StorageCopyEndpoint, storageRole.CopyHandler, "POST")
		server.AddHandler(settings.StorageMoveEndpoint, storageRole.MoveHandler, "POST")
		server.AddHandler(settings.StorageRenameEndpoint, storage.RenameHandler, "POST")
	}

	var routerRole *router.Router = nil
	if config.Roles.Router.Enable {
		routerRole, err = router.New(config, server)
		if err != nil {
			return nil, err
		}

		server.AddHandler(settings.RouterExplorerEndpoint, server.TokenizedHandler(routerRole.ExplorerHandler), "POST", "GET")
		server.AddHandler(settings.RouterFilesystemEndpoint, server.TokenizedHandler(routerRole.FilesystemHandler), "POST", "GET")
		server.AddHandler(settings.RouterDevicesEndpoint, server.TokenizedHandler(routerRole.DevicesHandler), "POST", "GET")

		server.AddHandler(settings.RouterTopologyEndpoint, server.TokenizedHandler(routerRole.GetTopologyHandler), "GET")
		server.AddHandler(settings.RouterNotifyEndpoint, server.TokenizedHandler(routerRole.NotifyHandler), "POST")
	}

	metricsRegistry := metrics.NewRegistry()
	server.AddHandler(settings.ServerMetricsEndpoint, metricsRegistry.Handler().ServeHTTP, "GET")

	return &Application{
		storageRole: storageRole,
		runnerRole:  nil,
		routerRole:  routerRole,
		config:      config,
		server:      server,
		lanListener: listener.New(config),
		lanObserver: observer.New(config),
	}, nil
}

func (app *Application) Run() error {
	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()
		app.server.Start()
	}()

	if err := app.server.WaitStart(); err != nil {
		return err
	}

	fmt.Printf("Platform is on %s\n", app.config.Roles.Runner.Platform)
	fmt.Printf("Server started on https://%s\n", app.config.Url)

	if app.config.Roles.Router.Enable {
		fmt.Printf("Authenticate on https://%s/auth\n", app.config.Url)
	}

	if app.storageRole != nil {
		if app.config.Url == "localhost"+settings.DefaultPort {
			if err := app.storageRole.NotifyRouter(app.config.Url); err != nil {
				logger.Error("[App] failed notify router: %s", err)
			}
		}
		go func() {
			addrs := app.lanObserver.Start()
			for {
				if err := app.storageRole.NotifyRouter(<-addrs + settings.DefaultPort); err != nil {
					logger.Error("[App] failed notify router: %s", err)
				}
			}
		}()
	}
	if app.routerRole != nil {
		app.lanListener.Start()
	}

	wg.Wait()

	return nil
}
