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
	connector, err := connector.NewConnector(config)
	if err != nil {
		return nil, err
	}

	server, err := server.New(config, connector)
	if err != nil {
		return nil, err
	}

	server.AddHandler(settings.ServerStatusEndpoint, server.StatusHandler, "GET")
	server.AddHandler(settings.ServerAuthEndpoint, server.AuthHandler, "GET")
	server.AddHandler(settings.ServerAuthTokenEndpoint, server.AuthTokenHandler, "POST")
	server.AddHandler(settings.ServerHomeEndpoint, server.HomeHandler, "GET")
	server.AddHandler(settings.ServerSettingsEndpoint, server.SettingsHandler, "GET")

	var storageRole *storage.Storage = nil
	if config.Roles.Storage.Enable {
		storageRole, err = storage.New(config, connector)
		if err != nil {
			return nil, fmt.Errorf("[App] failed create storage role: %s", err)
		}

		server.AddHandler(settings.StorageFilesystemEndpoint, server.TokenizedHandler(storageRole.FilesystemHandler), "POST", "GET")
		server.AddHandler(settings.StorageInsertEndpoint, server.TokenizedHandler(storageRole.InsertHandler), "POST")
		server.AddHandler(settings.StorageSelectEndpoint, server.TokenizedHandler(storageRole.SelectHandler), "POST")
		server.AddHandler(settings.StorageUpdateEndpoint, server.TokenizedHandler(storageRole.UpdateHandler), "POST")
		server.AddHandler(settings.StorageDeleteEndpoint, server.TokenizedHandler(storageRole.DeleteHandler), "POST")
		server.AddHandler(settings.StorageCopyEndpoint, server.TokenizedHandler(storageRole.CopyHandler), "POST")
		server.AddHandler(settings.StorageMoveEndpoint, server.TokenizedHandler(storageRole.MoveHandler), "POST")
		server.AddHandler(settings.StorageRenameEndpoint, server.TokenizedHandler(storage.RenameHandler), "POST")
	}

	var runnerRole *runner.Runner = nil
	if config.Roles.Runner.Enable {
		runnerRole, err = runner.New(config, connector)
		if err != nil {
			return nil, fmt.Errorf("[App] failed create runner role: %s", err)
		}

		server.AddHandler(settings.RunnerOpenEndpoint, server.TokenizedHandler(runnerRole.OpenFileHandler), "POST")

		server.AddHandler(settings.RunnerAppImageEndpoint, server.TokenizedHandler(runnerRole.AppImageHandler), "GET")
		server.AddHandler(settings.RunnerAppStreamEndpoint, server.TokenizedHandler(runnerRole.AppStreamHandler), "GET")
		server.AddHandler(settings.RunnerAppDirectStreamEndpoint, server.TokenizedHandler(runnerRole.AppDirectStreamHandler), "GET")
		server.AddHandler(settings.RunnerAppClickedEndpoint, server.TokenizedHandler(runnerRole.AppMouseClickedHandler), "POST")
	}

	var routerRole *router.Router = nil
	if config.Roles.Router.Enable {
		routerRole, err = router.New(config, server, connector)
		if err != nil {
			return nil, fmt.Errorf("[App] failed create router role: %s", err)
		}

		server.AddHandler(settings.RouterExplorerEndpoint, server.TokenizedHandler(routerRole.ExplorerHandler), "POST", "GET")
		server.AddHandler(settings.RouterFilesystemEndpoint, server.TokenizedHandler(routerRole.FilesystemHandler), "POST", "GET")
		server.AddHandler(settings.RouterDevicesEndpoint, server.TokenizedHandler(routerRole.DevicesHandler), "POST", "GET")

		server.AddHandler(settings.RouterTopologyEndpoint, server.TokenizedHandler(routerRole.GetTopologyHandler), "GET")
		server.AddHandler(settings.RouterNotifyEndpoint, server.TokenizedHandler(routerRole.NotifyHandler), "POST")

		server.AddHandler(settings.RouterOpenEndpoint, server.TokenizedHandler(routerRole.OpenFileHandler), "POST")

		server.AddHandler(settings.RouterFiltersGetEndpoint, server.TokenizedHandler(routerRole.FiltersGetHandler), "GET")
		server.AddHandler(settings.RouterFiltersAddEndpoint, server.TokenizedHandler(routerRole.FiltersAddHandler), "POST")
		server.AddHandler(settings.RouterFiltersRemoveEndpoint, server.TokenizedHandler(routerRole.FiltersRemoveHandler), "POST")
		server.AddHandler(settings.RouterFiltersSwapEndpoint, server.TokenizedHandler(routerRole.FiltersSwapHandler), "POST")
	}

	metricsRegistry := metrics.NewRegistry()
	server.AddHandler(settings.ServerMetricsEndpoint, metricsRegistry.Handler().ServeHTTP, "GET")

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
