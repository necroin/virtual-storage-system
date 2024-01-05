package settings

const (
	ServerStatusEndpoint    = "/status"
	ServerMetricsEndpoint   = "/metrics"
	ServerAuthEndpoint      = "/auth"
	ServerAuthTokenEndpoint = "/auth/token"
)

const (
	StorageFilesystemEndpoint = "/{token}/storage/filesystem"
	StorageInsertEndpoint     = "/{token}/storage/insert/{type}"
	StorageSelectEndpoint     = "/{token}/storage/select"
	StorageUpdateEndpoint     = "/{token}/storage/update"
	StorageDeleteEndpoint     = "/{token}/storage/delete"
	StorageCopyEndpoint       = "/{token}/storage/copy"
	StorageMoveEndpoint       = "/{token}/storage/move"
	StorageRenameEndpoint     = "/{token}/storage/rename"
)

const (
	RunnerOpenEndpoint   = "/runner/open"
	RunnerNotifyEndpoint = "/runner/notify"
)

const (
	RouterExplorerEndpoint   = "/{token}/router/explorer"
	RouterFilesystemEndpoint = "/{token}/router/filesystem"
	RouterDevicesEndpoint    = "/{token}/router/devices"
	RouterTopologyEndpoint   = "/{token}/router/topology"
	RouterNotifyEndpoint     = "/{token}/router/notify"
)
