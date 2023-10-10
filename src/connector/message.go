package connector

import (
	"time"
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

type FileInfo struct {
	ModTime time.Time `json:"mod_time"`
	Size    int64     `json:"size"`
}

type FilesystemDirectory struct {
	Directories map[string]FileInfo `json:"directories"`
	Files       map[string]FileInfo `json:"files"`
}
