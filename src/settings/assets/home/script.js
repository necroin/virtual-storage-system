function request(methood, url, data) {
    var req = new XMLHttpRequest();
    req.open(methood, url, false);
    req.send(data);
    return req.responseText
}

function GetTopology(url) {
    let response = request("GET", "https://"+url+"/router/topology")
    let data = JSON.parse(response)

    let storagesList = document.getElementById("storages")
    for (index in data.storages) {
        let storage = data.storages[index]

        let storageElement = document.createElement("div")
        storageElement.className = "list-item"
        let storageName = document.createElement("span")
        storageName.innerText = storage.hostname

        let storageUrl = document.createElement("span")
        storageUrl.innerText = storage.url

        storageElement.appendChild(storageName)
        storageElement.appendChild(storageUrl)

        storagesList.appendChild(storageElement)
    }
}