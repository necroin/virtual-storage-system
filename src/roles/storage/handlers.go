package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"strconv"
	"vss/src/buffer"
	"vss/src/settings"
	"vss/src/utils"
	"vss/src/utils/html"

	_ "embed"

	"github.com/gorilla/mux"
)

var (
	insertHandlers = map[string]func(string){
		"dir":      func(path string) { utils.CreateNewDirectory(path, "Новая папка") },
		"textFile": func(path string) { utils.CreateNewFile(path, "Текстовый документ.txt") },
	}
	copyHandlers = map[string]func(string){
		"file": func(path string) { buffer.SetFile(path) },
		"text": func(text string) {},
	}
)

func (storage *Storage) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	insertType := msgVars["type"]

	path, _ := ioutil.ReadAll(request.Body)

	insertHandler := insertHandlers[insertType]
	insertHandler(string(path))
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)
	utils.RemoveFile(string(msgPath))
}

func (storage *Storage) CopyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgVars := mux.Vars(request)
	copyType := msgVars["type"]

	path, _ := ioutil.ReadAll(request.Body)

	copyHandler := copyHandlers[copyType]
	copyHandler(string(path))
}

func (storage *Storage) PasteHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	fileSystemMessage := storage.CollectFileSystem(string(msgPath))
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func (storage *Storage) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	walkPath := "/Test"
	if len(msgPath) != 0 {
		walkPath = string(msgPath)
	}
	walkDirectory := storage.CollectFileSystem(walkPath)

	table_rows := html.NewTag("tbody").AddAttribute(html.NewAttribute("id", "filesystem-explorer-table-body"))
	rows_count := int64(0)

	directories := utils.GetMapKeys(walkDirectory.Directories)
	sort.Strings(directories)
	for _, directory := range directories {
		stat := walkDirectory.Directories[directory]
		row := html.NewTag("tr").AddAttribute(html.NewAttribute("tabindex", strconv.FormatInt(int64(rows_count), 10)))
		row_name := html.NewTag("td").AddElements(html.NewText("📁 " + directory))
		row_date := html.NewTag("td").AddElements(html.NewText(fmt.Sprintf("%v", stat.ModTime().Format("02.01.2006 15:04"))))
		row_type := html.NewTag("td").AddElements(html.NewText("Папка с файлами"))
		row_size := html.NewTag("td").AddElements(html.NewText(""))
		row.AddAttribute(
			html.NewAttribute("ondblclick", fmt.Sprintf("window.open('%s')", path.Join(walkPath, directory))),
			html.NewAttribute("name", directory),
		)
		row.AddElements(row_name, row_date, row_type, row_size)
		table_rows.AddElements(row)
		rows_count += 1
	}

	files := utils.GetMapKeys(walkDirectory.Files)
	sort.Strings(files)
	for _, file := range files {
		stat := walkDirectory.Files[file]
		row := html.NewTag("tr").AddAttribute(html.NewAttribute("tabindex", strconv.FormatInt(int64(rows_count), 10)))
		row_name := html.NewTag("td").AddElements(html.NewText(file))
		row_date := html.NewTag("td").AddElements(html.NewText(fmt.Sprintf("%v", stat.ModTime().Format("02.01.06 15:04"))))
		row_type := html.NewTag("td").AddElements(html.NewText("Файл"))
		row_size := html.NewTag("td").AddElements(html.NewText(fmt.Sprintf("%v байт", stat.Size())))
		row.AddAttribute(html.NewAttribute("name", file))
		row.AddElements(row_name, row_date, row_type, row_size)
		table_rows.AddElements(row)
		rows_count += 1
	}

	result := table_rows.InnerHTML()
	if len(msgPath) == 0 {
		style := html.NewTag("style").AddElements(html.NewText(settings.GetExplorerStyle())).AddAttribute(html.NewAttribute("type", "text/css"))
		script := html.NewScript(fmt.Sprintf(settings.GetExplorerScript(), "http://"+storage.url+settings.StorageMainEndpoint))
		result = fmt.Sprintf(
			settings.GetExplorerTemlate(),
			style.ToHTML(), script.ToHTML(),
			settings.ExplorerIconCreate, settings.ExplorerIconCut, settings.ExplorerIconCopy, settings.ExplorerIconPaste, settings.ExplorerIconDelete,
			settings.ExplorerIconArrowLeft,
			walkPath,
			table_rows.ToHTML(),
		)
	}

	responseWriter.Write([]byte(result))
}
