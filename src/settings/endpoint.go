package settings

const (
	ServerStatusEndpoint  = "/status"
	ServerMetricsEndpoint = "/metrics"
)

const (
	StorageFilesystemEndpoint = "/storage/filesystem"
	StorageViewEndpoint       = "/storage/view"
	StorageInsertEndpoint     = "/storage/insert"
	StorageSelectEndpoint     = "/storage/select"
	StorageUpdateEndpoint     = "/storage/update"
	StorageDeleteEndpoint     = "/storage/delete"
)

const (
	RunnerViewEndpoint     = "/runner/view"
	RunnerOpenEndpoint     = "/runner/open"
	RunnerTopologyEndpoint = "/runner/topology"
	RunnerNotifyEndpoint   = "/runner/notify"
)
const (
	RouterViewEndpoint     = "/router/view"
	RouterTopologyEndpoint = "/router/topology"
	RouterNotifyEndpoint   = "/router/notify"
)
