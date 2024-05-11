const filtersBlackListType = "Black list"
const filtersWhiteListType = "White list"

function request(method, url, data) {
    var req = new XMLHttpRequest();
    req.open(method, url, false);
    req.send(data);
    return req.responseText
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
}

function GetFilers(url) {
    let data = request("GET", 'https://' + url + "/router/filters/get")
    let filters = JSON.parse(data)
    return filters
}

function AddFiler(url, path) {
    if (path != "") {
        request("POST", 'https://' + url + "/router/filters/add", path)
        UpdateFilters(url, GetFilers(url))
    }
}

function RemoveFiler(url, path) {
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
            RemoveFiler(url, filter)
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

function GetDevices(url) {
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

function UpdateDevices() {

}