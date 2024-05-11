window.__context__ = {}

function Init() {
    let savedPath = window.localStorage.getItem("explorer-address-line")
    if (savedPath == null) {
        savedPath = "/"
    }
    SetCurrentPath(savedPath)
}

function remove_dropdown(event, event_id, dropdown_id) {
    if (!event.target.matches('#' + event_id)) {
        document.getElementById(dropdown_id).classList.remove('dropdown-show');
    }
}

window.onclick = function (event) {
    remove_dropdown(event, "bar-create-button", "create-options")
    remove_dropdown(event, "bar-options-button", "options")
}

function async_request(methood, url, data, callback) {
    var req = new XMLHttpRequest();
    req.onload = () => {
        if (req.readyState === XMLHttpRequest.DONE) {
            callback(req.responseText)
        }
    }
    req.open(methood, url, true);
    req.send(data);
}

function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function SetStorage(url) {
    window.__context__.storageUrl = url
}

function GetRequestUrl(routerUrl) {
    if (window.__context__.storageUrl != null) {
        return window.__context__.storageUrl
    }
    return routerUrl
}

function GetRequestRole() {
    if (window.__context__.storageUrl != null) {
        return "/storage"
    }
    return "/router"
}

function GetCurrentPath() {
    return document.getElementById("explorer-address-line").value
}

function SetCurrentPath(path) {
    document.getElementById("explorer-address-line").value = path
    window.localStorage.setItem("explorer-address-line", path)
}

function GetFilesystem(routerUrl, path) {
    let rowsCount = 0
    if (path != null) {
        SetCurrentPath(path)
    }
    path = GetCurrentPath()
    callback = (responseText) => {
        let filesystem = JSON.parse(responseText)
        let filters = GetFilers(routerUrl)

        let filesystemTable = document.getElementById("explorer-filesystem-content-body")
        filesystemTable.replaceChildren()
        filesystemTable.focusItem = null

        let directories = filesystem.directories
        for (let directory in directories) {
            if (!IsPassFilter(filters, path, directory)) {
                continue
            }

            let info = directories[directory]

            let tableRow = document.createElement("tr")
            tableRow.tabIndex = String(rowsCount)
            tableRow.__custom__ = {}
            tableRow.__custom__["name"] = directory

            let storageUrl = info["url"]
            let storageUrlParse = new URL("https://" + storageUrl)
            if (storageUrlParse.hostname == "localhost") {
                storageUrl = [location.host, storageUrlParse.pathname.split("/")[1]].join("/")
            }
            tableRow.__custom__["storageUrl"] = storageUrl

            let nameElement = document.createElement("td")
            let storageElement = document.createElement("td")
            let dateElement = document.createElement("td")
            let typeElement = document.createElement("td")
            let sizeElement = document.createElement("td")

            nameElement.innerText = "📁 " + directory
            dateElement.innerText = info["mod_time"]
            typeElement.innerText = "Directory"
            sizeElement.innerText = ""

            tableRow.appendChild(nameElement)
            tableRow.appendChild(storageElement)
            tableRow.appendChild(dateElement)
            tableRow.appendChild(typeElement)
            tableRow.appendChild(sizeElement)

            let openPath = "/" + directory
            if (path == "/") {
                openPath = directory
            }
            openPath = path + openPath

            tableRow.ondblclick = () => GetFilesystem(GetRequestUrl(routerUrl), openPath)
            tableRow.ontouchend = () => GetFilesystem(GetRequestUrl(routerUrl), openPath)

            filesystemTable.appendChild(tableRow)

            rowsCount = rowsCount + 1
        }

        let files = filesystem.files
        for (let file in files) {
            if (!IsPassFilter(filters, path, file)) {
                continue
            }

            let infos = files[file]
            for (infoIndex in infos) {
                let info = infos[infoIndex]

                let tableRow = document.createElement("tr")
                tableRow.tabIndex = String(rowsCount)
                tableRow.__custom__ = {}
                tableRow.__custom__["name"] = file
                tableRow.__custom__["platform"] = info["platform"]
                tableRow.__custom__["hostname"] = info["hostname"]

                let storageUrl = info["url"]
                let storageUrlParse = new URL("https://" + storageUrl)
                if (storageUrlParse.hostname == "localhost") {
                    storageUrl = [location.host, storageUrlParse.pathname.split("/")[1]].join("/")
                }
                tableRow.__custom__["storageUrl"] = storageUrl

                let nameElement = document.createElement("td")
                let storageElement = document.createElement("td")
                let dateElement = document.createElement("td")
                let typeElement = document.createElement("td")
                let sizeElement = document.createElement("td")

                nameElement.innerText = file
                storageElement.innerText = info["hostname"]
                dateElement.innerText = info["mod_time"]
                typeElement.innerText = "File"
                sizeElement.innerText = info["size"] + " KB"

                tableRow.appendChild(nameElement)
                tableRow.appendChild(storageElement)
                tableRow.appendChild(dateElement)
                tableRow.appendChild(typeElement)
                tableRow.appendChild(sizeElement)

                tableRow.ondblclick = () => OpenFile(routerUrl, tableRow)
                tableRow.ontouchend = () => OpenFile(routerUrl, tableRow)

                filesystemTable.appendChild(tableRow)

                rowsCount = rowsCount + 1
            }
        }
    }
    async_request("POST", "https://" + GetRequestUrl(routerUrl) + GetRequestRole() + "/filesystem", path, callback)
}

function GetDevices(routerUrl) {
    callback = (devicesResponse) => {
        let devices = JSON.parse(devicesResponse)

        let devicesList = document.getElementById("devices")
        devicesList.replaceChildren()
        let allDevicesElement = document.createElement("span")
        allDevicesElement.innerText = "All"
        allDevicesElement.onclick = () => {
            SetStorage(null)
            GetFilesystem(routerUrl)
        }
        devicesList.appendChild(allDevicesElement)

        let createOptions = document.getElementById("create-storage-select")

        for (let device in devices) {
            let deviceUrl = devices[device]
            let deviceUrlParse = new URL("https://" + deviceUrl)
            if (deviceUrlParse.hostname == "localhost") {
                deviceUrl = [location.host, deviceUrlParse.pathname.split("/")[1]].join("/")
            }

            let deviceElement = document.createElement("span")
            deviceElement.innerText = device
            deviceElement.onclick = () => {
                SetStorage(deviceUrl)
                GetFilesystem(routerUrl)
            }
            devicesList.appendChild(deviceElement)

            let createOptionDeviceElement = document.createElement("option")
            createOptionDeviceElement.innerText = device
            createOptionDeviceElement.__custom__ = {}
            createOptionDeviceElement.__custom__.storageUrl = deviceUrl
            createOptions.appendChild(createOptionDeviceElement)
        }
    }
    async_request("GET", "https://" + GetRequestUrl(routerUrl) + GetRequestRole() + "/devices", null, callback)
}

function Back(url) {
    let splitted_path = document.getElementById("explorer-address-line").value.split("/");
    let path = splitted_path.slice(0, splitted_path.length - 1).join("/");
    if (path == "") {
        path = "/"
    }
    GetFilesystem(url, path)
}

function SetFocusItem(table, item) {
    if (table.focusItem != null) {
        table.focusItem.style.backgroundColor = "var(--backgroud-color)"
    }
    table.focusItem = item
    table.focusItem.style.backgroundColor = "var(--elements-bg-color)"
}

function GetFocusItem() {
    return document.getElementById("explorer-filesystem-content-body").focusItem
}

function OpenCreateOptions() {
    document.getElementById("create-options").classList.toggle("dropdown-show");
}

function OpenOptions() {
    document.getElementById("options").classList.toggle("dropdown-show");
}

function OpenDialog(dialog, overlay) {
    document.getElementById(dialog).style.display = "flex";
    document.getElementById(overlay).style.display = "block";
}

function CloseDialog(dialog, overlay) {
    document.getElementById(dialog).style.display = "none";
    document.getElementById(overlay).style.display = "none";
}

function UpdateStatusBar(data) {
    document.getElementById("status-bar-progress").innerHTML = data.status
    document.getElementById("status-bar-text").innerText = data.text
}

function Create(type) {
    let createStorageSelect = document.getElementById("create-storage-select")
    let url = createStorageSelect.options[createStorageSelect.selectedIndex].__custom__.storageUrl
    let data = JSON.stringify(
        {
            "type": type,
            "path": GetCurrentPath(),
            "name": document.getElementById("create-dialog-name").value
        }
    );
    let response = request("POST", "https://" + url + "/storage/insert/" + type, data);
    CloseDialog('create-dialog', 'create-dialog-overlay')
    UpdateStatusBar(JSON.parse(response))
}

function Remove(routerUrl) {
    let focusItem = GetFocusItem()
    if (focusItem != null) {
        let focusItemName = focusItem.__custom__.name
        let focusItemStorageUrl = focusItem.__custom__.storageUrl
        let path = [GetCurrentPath(), focusItemName].join("/")
        if (GetCurrentPath() == "/") {
            path = "/" + focusItemName
        }
        async_request(
            "POST",
            "https://" + focusItemStorageUrl + "/storage/delete",
            path,
            (response) => {
                GetFilesystem(routerUrl)
                UpdateStatusBar(JSON.parse(response))
            }
        );
    }
}

function Rename(routerUrl) {
    let focusItem = GetFocusItem()
    if (focusItem != null) {
        let focusItemStorageUrl = focusItem.__custom__.storageUrl
        let oldName = focusItem.__custom__.name
        let newName = document.getElementById("rename-dialog-name").value
        let data = JSON.stringify({
            "path": GetCurrentPath(),
            "old_name": oldName,
            "new_name": newName
        })
        async_request(
            "POST",
            "https://" + focusItemStorageUrl + "/storage/rename",
            data,
            (response) => {
                GetFilesystem(routerUrl)
                UpdateStatusBar(JSON.parse(response))
            }
        );
        CloseDialog('rename-dialog', 'rename-dialog-overlay')
    }
}

function Cut() {
    let focusItem = GetFocusItem()
    window.__context__.paste = {
        path: GetCurrentPath(),
        name: focusItem.__custom__.name,
        url: focusItem.__custom__.storageUrl
    }
    window.__context__.paste_endpoint = "/storage/move"
}

function Copy() {
    let focusItem = GetFocusItem()
    window.__context__.paste = {
        path: GetCurrentPath(),
        name: focusItem.__custom__.name,
        url: focusItem.__custom__.storageUrl
    }
    window.__context__.paste_endpoint = "/storage/copy"
}

function Paste(routerUrl) {
    let pasteData = window.__context__.paste
    let pasteEndpoint = window.__context__.paste_endpoint
    async_request("POST",
        "https://" + GetRequestUrl(routerUrl) + pasteEndpoint,
        JSON.stringify({
            src_path: [pasteData.path, pasteData.name].join("/"),
            dst_path: [GetCurrentPath(), pasteData.name].join("/"),
            src_url: pasteData.url,
        }),
        (response) => {
            GetFilesystem(routerUrl)
            UpdateStatusBar(JSON.parse(response))
        }
    );
}

function OpenFile(routerUrl, item) {
    let platform = item.__custom__["platform"]
    let itemName = item.__custom__["name"]
    let path = [GetCurrentPath(), itemName].join("/")
    let srcUrl = item.__custom__["storageUrl"]
    let hostname = item.__custom__["hostname"]

    response = request("POST",
        "https://" + routerUrl + "/router/open",
        JSON.stringify({
            platform: platform,
            path: path,
            src_url: srcUrl,
            hostname: hostname,
        }),
    );

    let openResponse = JSON.parse(response)
    UpdateStatusBar(openResponse.status_bar)

    let runnerUrl = openResponse.runner_url
    let clientUrl = openResponse.client_url

    console.log(runnerUrl.split(":")[0])
    console.log(clientUrl)

    if (runnerUrl.split(":")[0] == clientUrl) {
        return
    }
    let pid = openResponse.pid
    window.open("https://" + runnerUrl + "/runner/stream/" + String(pid))
}

function GetFilers(routerUrl) {
    let data = request("GET", 'https://' + routerUrl + "/router/filters/get")
    let filters = JSON.parse(data)
    return filters
}

function IsPassFilter(filters, path, name) {
    if (filters.current_list == "Black list") {
        for (patternIndex in filters.black_list) {
            let pattern = filters.black_list[patternIndex]
            if (pattern == [path, name].join("/")) {
                return false
            }
            if (new RegExp(pattern).test(name)) {
                return false
            }
        }
        return true
    }
    if (filters.current_list == "White list") {
        for (patternIndex in filters.white_list) {
            let pattern = filters.white_list[patternIndex]
            if (pattern == [path, name].join("/")) {
                return true
            }
            if (new RegExp(pattern).test(name)) {
                return true
            }
        }
        return false
    }
    return false
}