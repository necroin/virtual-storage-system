package settings

const (
	ServerStatusEndpoint  = "/status"
	ServerMetricsEndpoint = "/metrics"
)

const (
	StorageMainEndpoint       = "/storage"
	StorageFilesystemEndpoint = "/storage/filesystem"
	StorageInsertEndpoint     = "/storage/insert"
	StorageSelectEndpoint     = "/storage/select"
	StorageUpdateEndpoint     = "/storage/update"
	StorageDeleteEndpoint     = "/storage/delete"
)

const (
	RunnerOpenEndpoint     = "/runner/open"
	RunnerTopologyEndpoint = "/runner/topology"
	RunnerNotifyEndpoint   = "/runner/notify"
)
const (
	RouterTopologyEndpoint = "/router/topology"
	RouterNotifyEndpoint   = "/router/notify"
)
