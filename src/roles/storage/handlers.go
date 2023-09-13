package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"vss/src/settings"
	"vss/src/utils/html"

	_ "embed"
)

func (storage *Storage) InsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) SelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) UpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) DeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	storage.db.InsertHandler(request.Body, responseWriter)
}

func (storage *Storage) FilesystemHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	fileSystemMessage := storage.CollectFileSystem(string(msgPath))
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func (storage *Storage) MainHandler(responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	walkPath := "/"
	if len(msgPath) != 0 {
		walkPath = string(msgPath)
	}
	walkDirectory := storage.CollectFileSystem(walkPath)

	table_rows := html.NewTag("tbody")
	rows_count := int64(0)

	for _, directory := range walkDirectory.Directories {
		row := html.NewTag("tr").AddAttribute(html.NewAttribute("tabindex", strconv.FormatInt(int64(rows_count), 10)))
		row_name := html.NewTag("td").AddElements(html.NewText("📁 " + directory))
		row_date := html.NewTag("td").AddElements(html.NewText("?"))
		row_type := html.NewTag("td").AddElements(html.NewText("Папка с файлами"))
		row_size := html.NewTag("td").AddElements(html.NewText("? байт"))
		row.AddAttribute(html.NewAttribute("ondblclick", fmt.Sprintf("window.open('%s')", path.Join(walkPath, directory))))
		row.AddElements(row_name, row_date, row_type, row_size)
		table_rows.AddElements(row)
		rows_count += 1
	}

	for _, file := range walkDirectory.Files {
		row := html.NewTag("tr").AddAttribute(html.NewAttribute("tabindex", strconv.FormatInt(int64(rows_count), 10)))
		row_name := html.NewTag("td").AddElements(html.NewText(file))
		row_date := html.NewTag("td").AddElements(html.NewText("?"))
		row_type := html.NewTag("td").AddElements(html.NewText("Файл"))
		row_size := html.NewTag("td").AddElements(html.NewText("? байт"))
		row.AddElements(row_name, row_date, row_type, row_size)
		table_rows.AddElements(row)
		rows_count += 1
	}

	style := html.NewTag("style").AddElements(html.NewText(settings.ExplorerStyle)).AddAttribute(html.NewAttribute("type", "text/css"))
	script := html.NewScript(fmt.Sprintf(settings.ExplorerScript, "http://"+storage.url+settings.StorageMainEndpoint))
	result := fmt.Sprintf(
		settings.ExplorerTemlate,
		style.ToHTML(), script.ToHTML(),
		settings.ExplorerIconCreate, settings.ExplorerIconCut, settings.ExplorerIconCopy, settings.ExplorerIconPaste, settings.ExplorerIconDelete,
		path.Join(walkPath, ".."),
		settings.ExplorerIconArrowLeft,
		walkPath,
		table_rows.ToHTML(),
	)

	responseWriter.Write([]byte(result))
}
