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
	Url     string    `json:"url"`
}

type FilesystemDirectory struct {
	Directories map[string]FileInfo `json:"directories"`
	Files       map[string]FileInfo `json:"files"`
}

type ClientRequest struct {
	Type string `json:"type"`
	Path string `json:"path"`
	Name string `json:"name"`
}

type StatusBarResponse struct {
	Status string `json:"status"`
	Text   string `json:"text"`
}
