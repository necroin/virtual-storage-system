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

type FilesystemDirectory struct {
	Directories map[string]*FilesystemDirectory `json:"directories"`
	Files       []string                        `json:"files"`
}

type FilesystemMessage FilesystemDirectory
