package roles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"sort"
	"strconv"
	"vss/src/connector"
	"vss/src/settings"
	"vss/src/utils"
	"vss/src/utils/html"
)

type Role interface {
	CollectFileSystem(walkPath string) connector.FilesystemDirectory
	GetUrl() string
	GetMainEndpoint() string
	GetHostnames() map[string]string
}

func FilesystemHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	fileSystemMessage := role.CollectFileSystem(string(msgPath))
	json.NewEncoder(responseWriter).Encode(fileSystemMessage)
}

func MainHandler(role Role, responseWriter http.ResponseWriter, request *http.Request) {
	msgPath, _ := ioutil.ReadAll(request.Body)

	walkPath := "/"
	if len(msgPath) != 0 {
		walkPath = string(msgPath)
	}
	walkDirectory := role.CollectFileSystem(walkPath)

	table_rows := html.NewTag("tbody").AddAttribute(html.NewAttribute("id", "filesystem-explorer-table-body"))
	rows_count := int64(0)

	directories := utils.GetMapKeys(walkDirectory.Directories)
	sort.Strings(directories)
	for _, directory := range directories {
		info := walkDirectory.Directories[directory]
		row := html.NewTag("tr").AddAttribute(html.NewAttribute("tabindex", strconv.FormatInt(int64(rows_count), 10)))
		row_name := html.NewTag("td").AddElements(html.NewText("📁 " + directory))
		row_date := html.NewTag("td").AddElements(html.NewText(fmt.Sprintf("%v", info.ModTime.Format("02.01.2006 15:04"))))
		row_type := html.NewTag("td").AddElements(html.NewText("Папка с файлами"))
		row_size := html.NewTag("td").AddElements(html.NewText(""))
		row.AddAttribute(
			html.NewAttribute("ondblclick", fmt.Sprintf("window.open('%s')", path.Join(walkPath, directory))),
			html.NewAttribute("name", directory),
			html.NewAttribute("custom_type", "dir"),
			html.NewAttribute("storage_url", info.Url),
		)
		row.AddElements(row_name, row_date, row_type, row_size)
		table_rows.AddElements(row)
		rows_count += 1
	}

	files := utils.GetMapKeys(walkDirectory.Files)
	sort.Strings(files)
	for _, file := range files {
		info := walkDirectory.Files[file]
		row := html.NewTag("tr").AddAttribute(html.NewAttribute("tabindex", strconv.FormatInt(int64(rows_count), 10)))
		row_name := html.NewTag("td").AddElements(html.NewText(file))
		row_date := html.NewTag("td").AddElements(html.NewText(fmt.Sprintf("%v", info.ModTime.Format("02.01.06 15:04"))))
		row_type := html.NewTag("td").AddElements(html.NewText("Файл"))
		row_size := html.NewTag("td").AddElements(html.NewText(fmt.Sprintf("%v байт", info.Size)))
		row.AddAttribute(
			html.NewAttribute("name", file),
			html.NewAttribute("custom_type", "file"),
			html.NewAttribute("storage_url", info.Url),
		)
		row.AddElements(row_name, row_date, row_type, row_size)
		table_rows.AddElements(row)
		rows_count += 1
	}

	result := table_rows.InnerHTML()
	if len(msgPath) == 0 {
		style := html.NewTag("style").AddElements(html.NewText(settings.GetExplorerStyle())).AddAttribute(html.NewAttribute("type", "text/css"))
		script := html.NewScript(fmt.Sprintf(settings.GetExplorerScript(), role.GetUrl()+role.GetMainEndpoint()))

		hostnames := role.GetHostnames()

		hostnamesList := html.NewDiv()
		hostnamesList.AddAttribute(html.NewAttribute("id", "devices"))
		allItem := html.NewSpan("Все")
		allItem.AddAttribute(html.NewAttribute("onclick", fmt.Sprintf("window.set_storage(null)")))
		hostnamesList.AddElements(allItem)

		hostnamesCreateSelect := html.NewTag("select")
		hostnamesCreateSelect.AddAttribute(html.NewAttribute("id", "create-storage-select"))

		for hostname, url := range hostnames {
			item := html.NewSpan(hostname)
			item.AddAttribute(html.NewAttribute("onclick", fmt.Sprintf("window.set_storage('%s/storage')", url)))
			hostnamesList.AddElements(item)

			selectOption := html.NewTag("option")
			selectOption.AddElements(html.NewText(hostname))
			selectOption.AddAttribute(html.NewAttribute("value", hostname))
			hostnamesCreateSelect.AddElements(selectOption)
		}

		result = fmt.Sprintf(
			settings.GetExplorerTemlate(),
			style.ToHTML(), script.ToHTML(),
			settings.ExplorerIconCreate,
			hostnamesCreateSelect.ToHTML(),
			settings.ExplorerIconCut,
			settings.ExplorerIconCopy,
			settings.ExplorerIconPaste,
			settings.ExplorerIconDelete,
			settings.ExplorerIconOptions,
			settings.ExplorerIconArrowLeft,
			walkPath,
			hostnamesList.ToHTML(),
			table_rows.ToHTML(),
			settings.ExplorerStatusBarSuccess,
		)
	}

	responseWriter.Write([]byte(result))
}
