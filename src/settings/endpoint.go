package settings

const (
	ServerStatusEndpoint    = "/status"
	ServerMetricsEndpoint   = "/metrics"
	ServerAuthEndpoint      = "/auth"
	ServerAuthTokenEndpoint = "/auth/token"
)

const (
	StorageMainEndpoint       = "/{token}/storage"
	StorageFilesystemEndpoint = "/{token}/storage/filesystem"
	StorageInsertEndpoint     = "/{token}/storage/insert/{type}"
	StorageSelectEndpoint     = "/{token}/storage/select"
	StorageUpdateEndpoint     = "/{token}/storage/update"
	StorageDeleteEndpoint     = "/{token}/storage/delete"
	StorageCopyEndpoint       = "/{token}/storage/copy/{type}"
	StoragePasteEndpoint      = "/{token}/storage/paste"
)

const (
	RunnerOpenEndpoint     = "/runner/open"
	RunnerTopologyEndpoint = "/runner/topology"
	RunnerNotifyEndpoint   = "/runner/notify"
)

const (
	RouterMainEndpoint     = "/{token}/router"
	RouterInsertEndpoint   = "/{token}/router/insert/{type}"
	RouterTopologyEndpoint = "/{token}/router/topology"
	RouterNotifyEndpoint   = "/{token}/router/notify"
)
