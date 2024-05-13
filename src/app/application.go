package app

import (
	"fmt"
	"sync"
	"vss/src/config"
	"vss/src/connector"
	"vss/src/lan/listener"
	"vss/src/lan/observer"
	"vss/src/logger"
	"vss/src/roles/router"
	"vss/src/roles/runner"
	"vss/src/roles/storage"
	"vss/src/server"
	"vss/src/settings"
	"vss/src/utils"

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

func New(config *config.Config) (*Application, error) {
	connector, err := connector.NewConnector(config.RootCAs)
	if err != nil {
		return nil, err
	}

	server, err := server.New(config, connector)
	if err != nil {
		return nil, err
	}

	server.AddHandlerFunc(settings.ServerStatusEndpoint, server.StatusHandler, "GET")
	server.AddHandlerFunc(settings.ServerAuthEndpoint, server.AuthHandler, "GET")
	server.AddHandlerFunc(settings.ServerAuthTokenEndpoint, server.AuthTokenHandler, "POST")
	server.AddHandlerFunc(settings.ServerHomeEndpoint, server.HomeHandler, "GET")
	server.AddHandlerFunc(settings.ServerSettingsEndpoint, server.SettingsHandler, "GET")

	var storageRole *storage.Storage = nil
	if config.Roles.Storage.Enable {
		storageRole, err = storage.New(config, connector)
		if err != nil {
			return nil, fmt.Errorf("[App] failed create storage role: %s", err)
		}

		server.AddHandlerFunc(settings.StorageFilesystemEndpoint, server.TokenizedHandler(storageRole.FilesystemHandler), "POST", "GET")
		server.AddHandlerFunc(settings.StorageInsertEndpoint, server.TokenizedHandler(storageRole.InsertHandler), "POST")
		server.AddHandlerFunc(settings.StorageSelectEndpoint, server.TokenizedHandler(storageRole.SelectHandler), "POST")
		server.AddHandlerFunc(settings.StorageUpdateEndpoint, server.TokenizedHandler(storageRole.UpdateHandler), "POST")
		server.AddHandlerFunc(settings.StorageDeleteEndpoint, server.TokenizedHandler(storageRole.DeleteHandler), "POST")
		server.AddHandlerFunc(settings.StorageCopyEndpoint, server.TokenizedHandler(storageRole.CopyHandler), "POST")
		server.AddHandlerFunc(settings.StorageMoveEndpoint, server.TokenizedHandler(storageRole.MoveHandler), "POST")
		server.AddHandlerFunc(settings.StorageRenameEndpoint, server.TokenizedHandler(storage.RenameHandler), "POST")
	}

	var runnerRole *runner.Runner = nil
	if config.Roles.Runner.Enable {
		runnerRole, err = runner.New(config, connector)
		if err != nil {
			return nil, fmt.Errorf("[App] failed create runner role: %s", err)
		}

		server.AddHandlerFunc(settings.RunnerOpenEndpoint, server.TokenizedHandler(runnerRole.OpenFileHandler), "POST")

		server.AddHandlerFunc(settings.RunnerAppImageEndpoint, server.TokenizedHandler(runnerRole.AppImageHandler), "GET")
		server.AddHandlerFunc(settings.RunnerAppStreamEndpoint, server.TokenizedHandler(runnerRole.AppStreamHandler), "GET")
		server.AddHandlerFunc(settings.RunnerAppDirectStreamEndpoint, server.TokenizedHandler(runnerRole.AppDirectStreamHandler), "GET")
		server.AddHandlerFunc(settings.RunnerAppMouseEventEndpoint, server.TokenizedHandler(runnerRole.AppMouseEventHandler), "POST")
		server.AddHandlerFunc(settings.RunnerAppKeyboardEventEndpoint, server.TokenizedHandler(runnerRole.AppKeyboardEventHandler), "POST")
	}

	var routerRole *router.Router = nil
	if config.Roles.Router.Enable {
		routerRole, err = router.New(config, server, connector)
		if err != nil {
			return nil, fmt.Errorf("[App] failed create router role: %s", err)
		}

		server.AddHandlerFunc(settings.RouterExplorerEndpoint, server.TokenizedHandler(routerRole.ExplorerHandler), "POST", "GET")
		server.AddHandlerFunc(settings.RouterFilesystemEndpoint, server.TokenizedHandler(routerRole.FilesystemHandler), "POST", "GET")
		server.AddHandlerFunc(settings.RouterDevicesEndpoint, server.TokenizedHandler(routerRole.DevicesHandler), "POST", "GET")

		server.AddHandlerFunc(settings.RouterTopologyEndpoint, server.TokenizedHandler(routerRole.GetTopologyHandler), "GET")
		server.AddHandlerFunc(settings.RouterNotifyEndpoint, server.TokenizedHandler(routerRole.NotifyHandler), "POST")

		server.AddHandlerFunc(settings.RouterOpenEndpoint, server.TokenizedHandler(routerRole.OpenFileHandler), "POST")

		server.AddHandlerFunc(settings.RouterFiltersGetEndpoint, server.TokenizedHandler(routerRole.FiltersGetHandler), "GET")
		server.AddHandlerFunc(settings.RouterFiltersAddEndpoint, server.TokenizedHandler(routerRole.FiltersAddHandler), "POST")
		server.AddHandlerFunc(settings.RouterFiltersRemoveEndpoint, server.TokenizedHandler(routerRole.FiltersRemoveHandler), "POST")
		server.AddHandlerFunc(settings.RouterFiltersSwapEndpoint, server.TokenizedHandler(routerRole.FiltersSwapHandler), "POST")

		server.AddHandlerFunc(settings.RouterReplicationGetEndpoint, server.TokenizedHandler(routerRole.ReplicationGetHandler), "GET")
		server.AddHandlerFunc(settings.RouterReplicationAddEndpoint, server.TokenizedHandler(routerRole.ReplicationAddHandler), "POST")
		server.AddHandlerFunc(settings.RouterReplicationRemoveEndpoint, server.TokenizedHandler(routerRole.ReplicationRemoveHandler), "POST")

		routerRole.AddReplicationTasks(config.Settings.Replication...)
	}

	metricsRegistry := metrics.NewRegistry()
	server.AddHandler(settings.ServerMetricsEndpoint, metricsRegistry.Handler(), "GET")

	return &Application{
		storageRole: storageRole,
		runnerRole:  runnerRole,
		routerRole:  routerRole,
		config:      config,
		server:      server,
		lanListener: listener.New(config),
		lanObserver: observer.New(config),
	}, nil
}

func (app *Application) Run() error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.server.Start()
	}()

	if err := app.server.WaitStart(); err != nil {
		return err
	}

	if app.storageRole != nil || app.runnerRole != nil {
		if app.config.Url == "localhost"+settings.DefaultPort {
			if err := utils.IfNotNil(
				app.storageRole,
				func(object *storage.Storage) error { return object.NotifyRouter(app.config.Url) },
			); err != nil {
				logger.Error("[App] failed notify router by storage: %s", err)
			}

			if err := utils.IfNotNil(
				app.runnerRole,
				func(object *runner.Runner) error { return object.NotifyRouter(app.config.Url) },
			); err != nil {
				logger.Error("[App] failed notify router by runner: %s", err)
			}
		}

		go func() {
			addrs, _ := app.lanListener.Start()

			for {
				notifyAddr := <-addrs

				if err := utils.IfNotNil(
					app.storageRole,
					func(object *storage.Storage) error { return object.NotifyRouter(notifyAddr) },
				); err != nil {
					logger.Error("[App] failed notify router by storage: %s", err)
				}

				if err := utils.IfNotNil(
					app.runnerRole,
					func(object *runner.Runner) error { return object.NotifyRouter(notifyAddr) },
				); err != nil {
					logger.Error("[App] failed notify router by runner: %s", err)
				}
			}
		}()
	}

	if app.routerRole != nil {
		app.lanObserver.Start()
		app.routerRole.StartReplicationTasks()
	}

	fmt.Printf("Platform is %s\n", app.config.Roles.Runner.Platform)
	fmt.Println("---")
	fmt.Printf("Server started on https://%s\n", app.config.Url)

	if app.config.Roles.Router.Enable {
		fmt.Println("---")
		fmt.Printf("Authenticate on https://%s/auth\n", app.config.Url)
		fmt.Println("---")
		fmt.Printf("Home on https://%s/%s/home\n", app.config.Url, app.config.User.Token)
		fmt.Println("---")
		fmt.Printf("Explore filesystem on https://%s/%s/router/explorer\n", app.config.Url, app.config.User.Token)
		fmt.Println("---")
	}

	wg.Wait()

	return nil
}
