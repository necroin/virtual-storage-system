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
    let filters = GetFilers(url)
    UpdateFilters(url, filters)
    UpdateDevices(url)
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
