const filtersBlackListType = "Black list"
const filtersWhiteListType = "White list"

function request(method, url, data) {
    var req = new XMLHttpRequest();
    req.open(method, url, false);
    req.send(data);
    return req.responseText
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

function Collapse(button, id) {
    if (button.innerText == "Collapse") {
        document.getElementById(id).style.display = "none"
        button.innerText = "Expand"
        return
    }

    if (button.innerText == "Expand") {
        document.getElementById(id).style.display = "block"
        button.innerText = "Collapse"
        return
    }
}

function GetSettings(url) {
    UpdateDevices(url)
    UpdateReplicationList(url)


    let filters = GetFilers(url)
    UpdateFilters(url, filters)
}

function GetFilers(url) {
    let data = request("GET", 'https://' + url + "/router/filters/get")
    let filters = JSON.parse(data)
    return filters
}

function AddFilter(url, path) {
    if (path != "") {
        request("POST", 'https://' + url + "/router/filters/add", path)
        UpdateFilters(url, GetFilers(url))
    }
}

function RemoveFilter(url, path) {
    request("POST", 'https://' + url + "/router/filters/remove", path)
    UpdateFilters(url, GetFilers(url))
}

function UpdateFilters(url, filters) {
    document.getElementById("settings-filters-list-button").innerText = filters.current_list
    let currentFilters = filters.black_list
    if (filters.current_list == filtersWhiteListType) {
        currentFilters = filters.white_list
    }

    let filtersList = document.getElementById("settings-filters-list")
    filtersList.replaceChildren()

    for (index in currentFilters) {
        let filter = currentFilters[index]

        let filterElement = document.createElement("div")
        filterElement.className = "list-item"

        let filterNameElement = document.createElement("span")
        filterNameElement.innerText = filter

        let filterButtonElement = document.createElement("button")
        filterButtonElement.innerText = "✖"
        filterButtonElement.onclick = () => {
            RemoveFilter(url, filter)
        }

        filterElement.appendChild(filterNameElement)
        filterElement.appendChild(filterButtonElement)
        filtersList.appendChild(filterElement)
    }
}

function SwapFiltersListType(url) {
    request("POST", 'https://' + url + "/router/filters/swap")
    UpdateFilters(url, GetFilers(url))
}

function UpdateDevices(url) {
    callback = (devicesResponse) => {
        let devices = JSON.parse(devicesResponse)

        let srcDevices = document.getElementById("settings-replication-src-devices")
        let dstDevices = document.getElementById("settings-replication-dst-devices")

        srcDevices.replaceChildren()
        dstDevices.replaceChildren()

        for (let device in devices) {
            let srcDeviceOption = document.createElement("option")
            srcDeviceOption.innerText = device
            srcDevices.appendChild(srcDeviceOption)

            let dstDeviceOption = document.createElement("option")
            dstDeviceOption.innerText = device
            dstDevices.appendChild(dstDeviceOption)
        }
    }
    async_request("GET", "https://" + url + "/router/devices", null, callback)
}

function AddReplication(url) {
    let srcDevicesSelect = document.getElementById("settings-replication-src-devices")
    let dstDevicesSelect = document.getElementById("settings-replication-dst-devices")
    let srcInput = document.getElementById("settings-replication-src-input")
    let dstInput = document.getElementById("settings-replication-dst-input")
    let cronSelect = document.getElementById("settings-replication-cron")

    let srcDevice = srcDevicesSelect.options[srcDevicesSelect.selectedIndex].innerText
    let dstDevice = dstDevicesSelect.options[dstDevicesSelect.selectedIndex].innerText
    let srcPath = srcInput.value
    let dstPath = dstInput.value
    let cron = cronSelect.options[cronSelect.selectedIndex].attributes.cron.value
    let cronName = cronSelect.options[cronSelect.selectedIndex].innerText

    let replication = {
        "src_hostname": srcDevice,
        "dst_hostname": dstDevice,
        "src_path": srcPath,
        "dst_path": dstPath,
        "cron": cron,
        "cron_name": cronName
    }

    console.log(replication)

    if (srcPath != "" && dstPath != "") {
        request("POST", 'https://' + url + "/router/replication/add", JSON.stringify(replication))
        UpdateReplicationList(url)
    }
}

function RemoveReplication(url, replication) {
    request("POST", 'https://' + url + "/router/replication/remove", JSON.stringify(replication))
    UpdateReplicationList(url)
}

function UpdateReplicationList(url) {
    callback = (response) => {
        let replicationList = JSON.parse(response)

        let list = document.getElementById("settings-replication-list")
        list.replaceChildren()

        for (index in replicationList) {
            let listItem = replicationList[index]

            let listItemElement = document.createElement("div")
            listItemElement.className = "list-item"

            let contentElement = document.createElement("div")
            contentElement.className = "vertical-layout"

            let periodElement = document.createElement("span")
            periodElement.innerText = "Period: " + listItem.cron_name 

            let devicesElement = document.createElement("span")
            devicesElement.innerText = "[ " + listItem.src_hostname + " ] → [ " + listItem.dst_hostname + " ] "

            let pathsElement = document.createElement("span")
            pathsElement.innerText = listItem.src_path + " → " + listItem.dst_path

            contentElement.appendChild(periodElement)
            contentElement.appendChild(devicesElement)
            contentElement.appendChild(pathsElement)


            let deleteButtonElement = document.createElement("button")
            deleteButtonElement.innerText = "✖"
            deleteButtonElement.onclick = () => {
                RemoveReplication(url, listItem)
            }

            listItemElement.appendChild(contentElement)
            listItemElement.appendChild(deleteButtonElement)
            list.appendChild(listItemElement)
        }
    }
    async_request("GET", "https://" + url + "/router/replication/get", null, callback)
}