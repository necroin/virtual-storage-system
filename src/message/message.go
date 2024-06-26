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
	Pid               int
}

type Notify struct {
	Type      string `json:"type"`
	Url       string `json:"url"`
	Hostname  string `json:"hostname"`
	Token     string `json:"token"`
	Platform  string `json:"platform"`
	Timestamp int64  `json:"timestamp"`
}

type Topology struct {
	Storages []Notify `json:"storages"`
	Runners  []Notify `json:"runners"`
}

type FileInfo struct {
	ModTime  string `json:"mod_time"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	Platform string `json:"platform"`
	Hostname string `json:"hostname"`
}

type FilesystemDirectory struct {
	Directories map[string]FileInfo   `json:"directories"`
	Files       map[string][]FileInfo `json:"files"`
}

type ClientAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InsertRequest struct {
	Type string `json:"type"`
	Path string `json:"path"`
	Name string `json:"name"`
}

type RenameRequest struct {
	Path    string `json:"path"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type CopyRequest struct {
	SrcPath string `json:"src_path"`
	DstPath string `json:"dst_path"`
	SrcUrl  string `json:"src_url"`
}

type MoveRequest CopyRequest

type OpenRequest struct {
	Platform string `json:"platform"`
	Path     string `json:"path"`
	SrcUrl   string `json:"src_url"`
	Hostname string `json:"hostname"`
}

type StatusBarResponse struct {
	Status string `json:"status"`
	Text   string `json:"text"`
}

type OpenResponse struct {
	Pid       int               `json:"pid"`
	RunnerUrl string            `json:"runner_url"`
	ClientUrl string            `json:"client_url"`
	Error     error             `json:"error"`
	StatusBar StatusBarResponse `json:"status_bar"`
}

type Coords struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type CoordsDelta Coords

type MouseEvent struct {
	Type   string      `json:"type"`
	Coords Coords      `json:"coords"`
	Scroll CoordsDelta `json:"scroll_delta"`
}

type KeyboardEvent struct {
	Type    string `json:"type"`
	Keycode uint16 `json:"keycode"`
}
