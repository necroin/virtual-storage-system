package settings

const (
	ServerStatusEndpoint    = "/status"
	ServerMetricsEndpoint   = "/metrics"
	ServerAuthEndpoint      = "/auth"
	ServerAuthTokenEndpoint = "/auth/token"
)

const (
	StorageMainEndpoint       = "/storage"
	StorageFilesystemEndpoint = "/storage/filesystem"
	StorageInsertEndpoint     = "/storage/insert/{type}"
	StorageSelectEndpoint     = "/storage/select"
	StorageUpdateEndpoint     = "/storage/update"
	StorageDeleteEndpoint     = "/storage/delete"
	StorageCopyEndpoint       = "/storage/copy/{type}"
	StoragePasteEndpoint      = "/storage/paste"
)

const (
	RunnerOpenEndpoint     = "/runner/open"
	RunnerTopologyEndpoint = "/runner/topology"
	RunnerNotifyEndpoint   = "/runner/notify"
)

const (
	RouterMainEndpoint          = "/router"
	RouterTokenizedMainEndpoint = "/{token}/router"
	RouterInsertEndpoint        = "/router/insert/{type}"
	RouterTopologyEndpoint      = "/router/topology"
	RouterNotifyEndpoint        = "/router/notify"
)
