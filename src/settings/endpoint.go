package settings

const (
	ServerStatusEndpoint    = "/status"
	ServerMetricsEndpoint   = "/metrics"
	ServerAuthEndpoint      = "/auth"
	ServerAuthTokenEndpoint = "/auth/token"
	ServerHomeEndpoint      = "/{token}/home"
	ServerSettingsEndpoint  = "/{token}/settings"
)

const (
	StorageFilesystemEndpoint = "/{token}/storage/filesystem"
	StorageInsertEndpoint     = "/{token}/storage/insert/{type}"
	StorageSelectEndpoint     = "/{token}/storage/select"
	StorageUpdateEndpoint     = "/{token}/storage/update"
	StorageDeleteEndpoint     = "/{token}/storage/delete"
	StorageCopyEndpoint       = "/{token}/storage/copy/{type}"
	StorageMoveEndpoint       = "/{token}/storage/move/{type}"
	StorageRenameEndpoint     = "/{token}/storage/rename"
)

const (
	RunnerOpenEndpoint   = "/{token}/runner/open"
	RunnerNotifyEndpoint = "/{token}/runner/notify"
)

const (
	RouterExplorerEndpoint      = "/{token}/router/explorer"
	RouterFilesystemEndpoint    = "/{token}/router/filesystem"
	RouterDevicesEndpoint       = "/{token}/router/devices"
	RouterTopologyEndpoint      = "/{token}/router/topology"
	RouterNotifyEndpoint        = "/{token}/router/notify"
	RouterOpenEndpoint          = "/{token}/router/open"
	RouterFiltersGetEndpoint    = "/{token}/router/filters/get"
	RouterFiltersAddEndpoint    = "/{token}/router/filters/add"
	RouterFiltersRemoveEndpoint = "/{token}/router/filters/remove"
)
