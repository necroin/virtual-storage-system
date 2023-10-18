package app

import (
	"fmt"
	"sync"
	"vss/src/config"
	"vss/src/lan"
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
	lanListener *lan.Listener
	lanObserver *lan.Observer
}

func New() (*Application, error) {
	config, err := config.Load(settings.ConfigPath)
	if err != nil {
		return nil, err
	}

	server := server.New(config)

	server.AddHandler(settings.ServerStatusEndpoint, server.StatusHandler, "GET")
	server.AddHandler(settings.ServerAuthEndpoint, server.AuthHandler, "GET")
	server.AddHandler(settings.ServerAuthTokenEndpoint, server.AuthTokenHandler, "POST")

	var storageRole *storage.Storage = nil
	if config.Roles.Storage.Enable {
		storageRole, err = storage.New(config, "storage.db")
		if err != nil {
			return nil, err
		}

		server.AddHandler(settings.StorageMainEndpoint, storageRole.MainHandler, "POST", "GET")
		server.AddHandler(settings.StorageFilesystemEndpoint, storageRole.FilesystemHandler, "GET")
		server.AddHandler(settings.StorageInsertEndpoint, storageRole.InsertHandler, "POST")
		server.AddHandler(settings.StorageSelectEndpoint, storageRole.SelectHandler, "POST")
		server.AddHandler(settings.StorageUpdateEndpoint, storageRole.UpdateHandler, "POST")
		server.AddHandler(settings.StorageDeleteEndpoint, storageRole.DeleteHandler, "POST")
		server.AddHandler(settings.StorageCopyEndpoint, storageRole.CopyHandler, "POST")
		server.AddHandler(settings.StoragePasteEndpoint, storageRole.PasteHandler, "POST")
	}

	var routerRole *router.Router = nil
	if config.Roles.Router.Enable {
		routerRole, err = router.New(config)
		if err != nil {
			return nil, err
		}

		server.AddHandler(settings.RouterMainEndpoint, server.TokenizedHandler(routerRole.MainHandler), "POST", "GET")
		server.AddHandler(settings.RouterTopologyEndpoint, server.TokenizedHandler(routerRole.GetTopologyHandler), "GET")
		server.AddHandler(settings.RouterNotifyEndpoint, server.TokenizedHandler(routerRole.NotifyHandler), "POST")
		server.AddHandler(settings.RouterInsertEndpoint, server.TokenizedHandler(routerRole.InsertHandler), "POST")
	}

	metricsRegistry := metrics.NewRegistry()
	server.AddHandler(settings.ServerMetricsEndpoint, metricsRegistry.Handler().ServeHTTP, "GET")

	return &Application{
		storageRole: storageRole,
		runnerRole:  nil,
		routerRole:  routerRole,
		config:      config,
		server:      server,
		lanListener: lan.NewListener(config),
		lanObserver: lan.NewObserver(config),
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
	fmt.Println("Server started.")

	if app.storageRole != nil {
		go func() {
			addrs := app.lanObserver.Start()
			for {
				if err := app.storageRole.NotifyRouter(<-addrs + settings.DefaultPort); err != nil {
					fmt.Println(err)
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
