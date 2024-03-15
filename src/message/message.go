package message

const (
	NotifyMessageStorageType = "storage"
	NotifyMessageRunnerType  = "runner"
)

type PageInfo struct {
	Url               string
	Token             string
	Style             string
	Script            string
	IconCreate        string
	IconCut           string
	IconCopy          string
	IconPaste         string
	IconDelete        string
	IconOptions       string
	IconArrowLeft     string
	StatusBarIcon     string
	BarHomeIcon       string
	BarFilesystemIcon string
	BarSettingsIcon   string
}

type NotifyMessage struct {
	Type      string `json:"type"`
	Url       string `json:"url"`
	Hostname  string `json:"hostname"`
	Token     string `json:"token"`
	Platform  string `json:"platform"`
	Timestamp int64  `json:"timestamp"`
}

type TopologyMessage struct {
	Storages []NotifyMessage `json:"storages"`
	Runners  []NotifyMessage `json:"runners"`
}

type FileInfo struct {
	ModTime  string `json:"mod_time"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	Platform string `json:"platform"`
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

type ClientAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RenameRequest struct {
	Path    string `json:"path"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type CopyRequest struct {
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	SrcUrl  string `json:"src_url"`
}

type OpenRequest struct {
	Platform string `json:"platform"`
	Path     string `json:"path"`
	SrcUrl   string `json:"src_url"`
	Type     string `json:"type"`
}

type OpenResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}
