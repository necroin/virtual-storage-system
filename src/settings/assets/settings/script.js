const filtersBlackListStorageName = "vss-settings-filters-black-list"
const filtersWhiteListStorageName = "vss-settings-filters-white-list"
const filtersCurrentListStorageName = "vss-settings-filters-current-list"

const filtersBlackListType = "Чёрный список"
const filtersWhiteListType = "Белый список"
const filtersCurrentListType = "Current List"

function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function StringToArray(value) {
    return !value ? [] : value.split(',')
}

function Init() {
    let filters = GetFilers()
    SaveFilters(filters)
    document.getElementById("settings-filters-list-button").innerText = filters[filtersCurrentListType]
}

function GetSettings(url) {
    let filters = GetFilers()
    UpdateFiltersList(filters)
}

function GetFilers() {
    let filters = {
        [filtersBlackListType]: [],
        [filtersWhiteListType]: [],
        [filtersCurrentListType]: filtersBlackListType
    }

    let blackList = window.localStorage.getItem(filtersBlackListStorageName)
    if (blackList != null) {
        filters[filtersBlackListType] = StringToArray(blackList)
    }

    let whiteList = window.localStorage.getItem(filtersWhiteListStorageName)
    if (whiteList != null) {
        filters[filtersWhiteListType] = StringToArray(whiteList)
    }

    let currentList = window.localStorage.getItem(filtersCurrentListStorageName)
    if (currentList != null) {
        filters[filtersCurrentListType] = currentList
    }
    console.log(filters)

    return filters
}

function SaveFilters(filters) {
    window.localStorage.setItem(filtersBlackListStorageName, filters[filtersBlackListType])
    window.localStorage.setItem(filtersWhiteListStorageName, filters[filtersWhiteListType])
    window.localStorage.setItem(filtersCurrentListStorageName, filters[filtersCurrentListType])
}


function AddFiler(path) {
    if (path != "") {
        let filters = GetFilers()
        let currentFiltersType = filters[filtersCurrentListType]
        filters[currentFiltersType].push(path)
        SaveFilters(filters)
        UpdateFiltersList(filters)
        document.getElementById('settings-filters-input').value = ""
    }
}

function RemoveFiler(path) {
    let filters = GetFilers()
    let currentFiltersType = filters[filtersCurrentListType]
    filters[currentFiltersType] = filters[currentFiltersType].filter((element) => { return element != path })
    SaveFilters(filters)
    UpdateFiltersList(filters)
}

function UpdateFiltersList(filters) {
    let currentFiltersType = filters[filtersCurrentListType]
    let currentFilters = filters[currentFiltersType]

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
            RemoveFiler(filter)
        }

        filterElement.appendChild(filterNameElement)
        filterElement.appendChild(filterButtonElement)
        filtersList.appendChild(filterElement)
    }
}

function SwapFiltersListType(button) {
    if (button.innerText == filtersBlackListType) {
        button.innerText = filtersWhiteListType
    } else {
        button.innerText = filtersBlackListType
    }
    window.localStorage.setItem(filtersCurrentListStorageName, button.innerText)
    UpdateFiltersList(GetFilers())
}