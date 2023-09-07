package connector

const (
	NotifyMessageStorageType = "storage"
	NotifyMessageRunnerType  = "runner"
)

type NotifyMessage struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type TopologyMessage struct {
	Storages []string `json:"storages"`
	Runners  []string `json:"runners"`
}
