package connector

import (
	"io/fs"
)

const (
	NotifyMessageStorageType = "storage"
	NotifyMessageRunnerType  = "runner"
)

type NotifyMessage struct {
	Type     string `json:"type"`
	Url      string `json:"url"`
	Hostname string `json:"hostname"`
}

type TopologyMessage struct {
	Storages []string `json:"storages"`
	Runners  []string `json:"runners"`
}

type FilesystemDirectory struct {
	Directories map[string]fs.FileInfo `json:"directories"`
	Files       map[string]fs.FileInfo `json:"files"`
}
